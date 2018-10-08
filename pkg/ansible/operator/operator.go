// Copyright 2018 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package operator

import (
	"log"
	"math/rand"
	"time"

	"github.com/operator-framework/operator-sdk/pkg/ansible/controller"
	"github.com/operator-framework/operator-sdk/pkg/ansible/runner"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/sirupsen/logrus"
)

func RunSDK(done chan error, mgr manager.Manager, namespace string) {
	if namespace == "" {
		namespace = "default"
	}
	watches, err := runner.NewFromWatches("./watches.yaml")
	if err != nil {
		logrus.Error("Failed to get watches")
		done <- err
		return
	}
	rand.Seed(time.Now().Unix())
	c := signals.SetupSignalHandler()

	for gvk, runner := range watches {
		controller.Add(mgr, controller.Options{
			GVK:         gvk,
			Namespace:   namespace,
			Runner:      runner,
			StopChannel: c,
		})
	}
	log.Fatal(mgr.Start(c))
	done <- nil
}
