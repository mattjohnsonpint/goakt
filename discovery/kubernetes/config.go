/*
 * MIT License
 *
 * Copyright (c) 2022-2025  Arsene Tochemey Gandote
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package kubernetes

import (
	"github.com/tochemey/goakt/v3/internal/validation"
)

// Config defines the kubernetes discovery configuration
type Config struct {
	// Namespace specifies the kubernetes namespace
	Namespace string
	// ApplicationName specifies the application name
	// Deprecated: this field is no longer used and will be removed in a future release. Instead, use the PodLabels field.
	ApplicationName string
	// ActorSystemName specifies the given actor system name
	// Deprecated: this field is no longer used and will be removed in a future release. Instead, use the PodLabels field.
	ActorSystemName string
	// DiscoveryPortName specifies the gossip port name
	DiscoveryPortName string
	// RemotingPortName specifies the remoting port name
	RemotingPortName string
	// PeersPortName specifies the cluster port name
	PeersPortName string

	// PodLabels specifies the pod labels
	PodLabels map[string]string
}

// Validate checks whether the given discovery configuration is valid
func (x Config) Validate() error {
	return validation.New(validation.FailFast()).
		AddValidator(validation.NewEmptyStringValidator("Namespace", x.Namespace)).
		AddValidator(validation.NewEmptyStringValidator("DiscoveryPortName", x.DiscoveryPortName)).
		AddValidator(validation.NewEmptyStringValidator("PeersPortName", x.PeersPortName)).
		AddValidator(validation.NewEmptyStringValidator("RemotingPortName", x.RemotingPortName)).
		Validate()
}
