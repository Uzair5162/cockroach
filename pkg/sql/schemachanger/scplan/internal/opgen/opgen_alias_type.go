// Copyright 2022 The Cockroach Authors.
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
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scop"
	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/scpb"
)

func init() {
	opRegistry.register((*scpb.AliasType)(nil),
		toPublic(
			scpb.Status_ABSENT,
			equiv(scpb.Status_DROPPED),
			to(scpb.Status_TXN_DROPPED,
				emit(func(this *scpb.AliasType) *scop.NotImplemented {
					return notImplemented(this)
				}),
			),
			to(scpb.Status_PUBLIC,
				emit(func(this *scpb.AliasType) *scop.MarkDescriptorAsPublic {
					return &scop.MarkDescriptorAsPublic{
						DescriptorID: this.TypeID,
					}
				}),
			),
		),
		toAbsent(
			scpb.Status_PUBLIC,
			to(scpb.Status_TXN_DROPPED,
				emit(func(this *scpb.AliasType, md *opGenContext) *scop.MarkDescriptorAsSyntheticallyDropped {
					return &scop.MarkDescriptorAsSyntheticallyDropped{
						DescriptorID: this.TypeID,
					}
				}),
			),
			to(scpb.Status_DROPPED,
				revertible(false),
				emit(func(this *scpb.AliasType) *scop.MarkDescriptorAsDropped {
					return &scop.MarkDescriptorAsDropped{
						DescriptorID: this.TypeID,
					}
				}),
			),
			to(scpb.Status_ABSENT,
				emit(func(this *scpb.AliasType, md *opGenContext) *scop.LogEvent {
					return newLogEventOp(this, md)
				}),
				emit(func(this *scpb.AliasType) *scop.DeleteDescriptor {
					return &scop.DeleteDescriptor{
						DescriptorID: this.TypeID,
					}
				}),
			),
		),
	)
}
