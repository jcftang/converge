// Copyright © 2016 Asteris, LLC
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

package image

import (
	"time"

	"github.com/asteris-llc/converge/resource"
	"github.com/asteris-llc/converge/resource/docker"
)

// Preparer for docker images
//
// Image is responsible for pulling Docker images. It assumes that there is
// already a Docker daemon running on the system.
type Preparer struct {
	// name of the image to pull
	Name string `hcl:"name"`

	// tag of the image to pull
	Tag string `hcl:"tag"`

	// the amount of time the pull will wait before halting forcefully. The
	// format is Go's duraction string. A duration string is a possibly signed
	// sequence of decimal numbers, each with optional fraction and a unit
	// suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns",
	// "us" (or "µs"), "ms", "s", "m", "h".
	Timeout string `hcl:"timeout" doc_type:"duration_string"`
}

// Prepare a new docker image
func (p *Preparer) Prepare(render resource.Renderer) (resource.Task, error) {
	name, err := render.Render("name", p.Name)
	if err != nil {
		return nil, err
	}

	tag, err := render.Render("tag", p.Tag)
	if err != nil {
		return nil, err
	}

	timeout, err := render.Render("timeout", p.Timeout)
	if err != nil {
		return nil, err
	}

	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		return nil, err
	}

	if timeout != "" {
		duration, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, err
		}
		dockerClient.PullInactivityTimeout = duration
	}

	image := &Image{
		Name: name,
		Tag:  tag,
	}
	image.SetClient(dockerClient)
	return image, nil
}
