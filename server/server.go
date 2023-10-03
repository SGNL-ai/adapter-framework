// Copyright 2023 SGNL.ai, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/server/internal"
)

// New returns an AdapterServer that wraps the given high-level
// Adapter implementation with the Tokens field populated from the file
// which name is configured in the AUTH_TOKENS_PATH environment variable.
// The stop channel returned is used to signal when the file watcher should
// be closed and stop watching for file changes.
func New[Config any](
	adapters map[string]framework.Adapter[Config],
) (api_adapter_v1.AdapterServer, chan struct{}) {
	authTokensPath, exists := os.LookupEnv("AUTH_TOKENS_PATH")
	if !exists {
		panic("AUTH_TOKENS_PATH environment variable not set")
	}

	return newWithAuthTokensPath(authTokensPath, adapters)

}

func newWithAuthTokensPath[Config any](
	authTokensPath string,
	adapters map[string]framework.Adapter[Config],
) (api_adapter_v1.AdapterServer, chan struct{}) {
	stop := make(chan struct{})

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(fmt.Sprintf("failed to create file watcher: %s", err.Error()))
	}

	if err = watcher.Add(authTokensPath); err != nil {
		panic(fmt.Sprintf("failed to add path to file watcher: %v", err))
	}

	server := &internal.Server[Config]{
		Adapters: adapters,
		Tokens:   getTokensFromPath(authTokensPath),
	}

	go func(s *internal.Server[Config]) {
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					// Channel was closed
					panic("file watcher channel closed")
				}

				s.TokensMutex.Lock()
				s.Tokens = getTokensFromPath(authTokensPath)
				s.TokensMutex.Unlock()
			case err, ok := <-watcher.Errors:
				if !ok {
					// Channel was closed
					panic("file watcher channel closed")
				}

				// An error will be thrown in the event there are too many events, too small of a buffer,
				// etc. This indicates the watcher may no longer be functioning correctly, so we'll panic.
				watcher.Close()
				panic(fmt.Sprintf("file watcher error: %v", err))
			case <-stop:
				watcher.Close()

				return
			}
		}
	}(server)

	return server, stop
}

// getTokensFromPath reads and parses the JSON encoded data located in the file at the given path.
// If the file does not exist or the file does not contain valid JSON, an empty slice is returned.
func getTokensFromPath(path string) []string {
	jsonValidTokens, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var validTokens *[]string

	if err := json.Unmarshal(jsonValidTokens, &validTokens); err != nil || validTokens == nil {
		return nil
	}

	return *validTokens
}
