// Copyright 2022 Red Hat, Inc.
// All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package shipper

import "github.com/spf13/cobra"

var root = &cobra.Command{
	Version: "0.0.1",
	Use:     "beatshipper",
	Short:   "beatshipper - Send data from files using beats",
	Long: `Sends data based on paths that will be exploded using GLOB
with the possibility of passing GNU Zip files also.`,
}

func Execute() error {
	return root.Execute()
}
