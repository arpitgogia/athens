package mongo

import (
	"context"

	"github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/observ"
	"github.com/gomods/athens/pkg/storage"
	"github.com/mongodb/mongo-go-driver/bson"
)

// List lists all versions of a module
func (s *ModuleStore) List(ctx context.Context, module string) ([]string, error) {
	const op errors.Op = "mongo.List"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()
	c := s.client.DB(s.db).C(s.coll)
	fields := bson.M{"version": 1}
	compositeQ := bson.D{bson.M{"module": module}, fields}
	cur, err := c.Find(compositeQ)
	if err != nil {
		kind := errors.KindUnexpected
		if err == mgo.ErrNotFound {
			kind = errors.KindNotFound
		}
		return nil, errors.E(op, kind, errors.M(module), err)
	}
	result := make([]storage.Module, 0)
	for cursor.Next() {
		var module storage.Module
		bytes, err := cursor.DecodeBytes()
		if err == nil {
			err = bson.Unmarshal(bytes, &module)
			if err == nil {
				result = append(result, module)
			}
		}
	}

	versions := make([]string, len(result))
	for i, r := range result {
		versions[i] = r.Version
	}

	return versions, nil
}
