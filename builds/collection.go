/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builds

import (
	"strings"

	"github.com/r3labs/diff"
)

type UniqueCollection []string

func (uc *UniqueCollection) add(item string) {
	for i := 0; i < len((*uc)); i++ {
		if (*uc)[i] == item {
			return
		}
	}

	(*uc) = append((*uc), item)
}

func componentsByType(t string, cl diff.Changelog) []string {
	var uc UniqueCollection
	for _, c := range cl.Filter([]string{t + "::*"}) {
		name := strings.TrimPrefix(c.Path[0], t+"::")
		uc.add(name)
	}

	return uc
}
