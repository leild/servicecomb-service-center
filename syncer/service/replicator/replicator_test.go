package replicator

import (
	"context"
	"testing"

	v1sync "github.com/apache/servicecomb-service-center/syncer/api/v1"
	"github.com/apache/servicecomb-service-center/syncer/service/replicator/resource"

	"github.com/stretchr/testify/assert"
)

func Test_replicatorManager_Persist(t *testing.T) {
	ctx := context.TODO()
	el := &v1sync.EventList{
		Events: nil,
	}
	r := manager.Persist(ctx, el)
	assert.Equal(t, 0, len(r))

	resource.RegisterResources("fork", func(event *v1sync.Event) resource.Resource {
		return &mockResources{
			loadCurrentResourceResult: nil,
			needOperateResult:         nil,
			operateResult:             resource.SuccessResult(),
		}
	})

	el = &v1sync.EventList{
		Events: []*v1sync.Event{
			{
				Id:        "xxx1",
				Action:    "",
				Subject:   "fork",
				Opts:      nil,
				Value:     nil,
				Timestamp: v1sync.Timestamp(),
			},
		},
	}

	r = manager.Persist(ctx, el)
	if assert.Equal(t, 1, len(r)) {
		assert.Equal(t, resource.SuccessResult().WithEventID("xxx1"), r[0])
	}

	el = &v1sync.EventList{
		Events: []*v1sync.Event{
			{
				Id:        "xxx1",
				Action:    "",
				Subject:   "not exist",
				Opts:      nil,
				Value:     nil,
				Timestamp: v1sync.Timestamp(),
			},
		},
	}

	r = manager.Persist(ctx, el)
	if assert.Equal(t, 1, len(r)) {
		assert.Equal(t, resource.Skip, r[0].Status)
	}
}

type mockResources struct {
	loadCurrentResourceResult *resource.Result
	needOperateResult         *resource.Result
	operateResult             *resource.Result
}

func (f *mockResources) LoadCurrentResource(_ context.Context) *resource.Result {
	return f.loadCurrentResourceResult
}

func (f *mockResources) NeedOperate(_ context.Context) *resource.Result {
	return f.needOperateResult
}

func (f *mockResources) Operate(_ context.Context) *resource.Result {
	return f.operateResult
}

func (f mockResources) FailHandle(_ context.Context, _ int32) (*v1sync.Event, error) {
	return nil, nil
}

func (f mockResources) CanDrop() bool {
	return true
}

func Test_pageEvents(t *testing.T) {
	t.Run("no page case", func(t *testing.T) {
		source := &v1sync.EventList{
			Events: []*v1sync.Event{
				{
					Value: []byte("hello"),
				},
			},
		}
		result := pageEvents(source, maxSize)
		if assert.Equal(t, 1, len(result)) {
			assert.Equal(t, []byte("hello"), result[0].Events[0].Value)
		}
	})

	t.Run("page case", func(t *testing.T) {
		source := &v1sync.EventList{
			Events: []*v1sync.Event{
				{
					Value: []byte("11"),
				},
				{
					Value: []byte("2"),
				},
			},
		}
		result := pageEvents(source, 3)
		if assert.Equal(t, 2, len(result)) {
			assert.Equal(t, []byte("11"), result[0].Events[0].Value)
			assert.Equal(t, []byte("2"), result[1].Events[0].Value)
		}
	})
}
