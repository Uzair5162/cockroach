// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package opgen

import (
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/tabledesc"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scop"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
)

func init() {
	opRegistry.register(
		(*scpb.IndexName)(nil),
		add(
			to(scpb.Status_PUBLIC,
				minPhase(scop.PreCommitPhase),
				emit(func(this *scpb.IndexName) scop.Op {
					return &scop.SetIndexName{
						TableID: this.TableID,
						IndexID: this.IndexID,
						Name:    this.Name,
					}
				}),
			),
		),
		drop(
			to(scpb.Status_ABSENT,
				minPhase(scop.PreCommitPhase),
				emit(func(this *scpb.IndexName) scop.Op {
					return &scop.SetIndexName{
						TableID: this.TableID,
						IndexID: this.IndexID,
						Name:    tabledesc.IndexNamePlaceholder(this.IndexID),
					}
				}),
			),
		),
	)

}
