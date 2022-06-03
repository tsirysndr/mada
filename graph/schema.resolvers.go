package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/tsirysndr/mada/graph/generated"
	"github.com/tsirysndr/mada/graph/model"
	"github.com/tsirysndr/mada/types"
)

func (r *mutationResolver) Empty(ctx context.Context) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Commune(ctx context.Context, id string) (*model.Commune, error) {
	res, err := r.CommuneService.ShowCommune(id)
	fmt.Println(res, err)
	return &model.Commune{}, nil
}

func (r *queryResolver) Communes(ctx context.Context, after *string, size *int) (*model.CommuneList, error) {
	res, err := r.CommuneService.List(0, 10)
	fmt.Println(res, err)
	return &model.CommuneList{}, nil
}

func (r *queryResolver) CountCommunes(ctx context.Context) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) District(ctx context.Context, id string) (*model.District, error) {
	r.DistrictService.ShowDistrict(id)
	return &model.District{}, nil
}

func (r *queryResolver) Districts(ctx context.Context, after *string, size *int) (*model.DistrictList, error) {
	r.DistrictService.List(0, 100)
	return &model.DistrictList{}, nil
}

func (r *queryResolver) CountDistricts(ctx context.Context) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Fokontany(ctx context.Context, id string) (*model.Fokontany, error) {
	r.FokontanyService.ShowFokontany(id)
	return &model.Fokontany{}, nil
}

func (r *queryResolver) AllFokontany(ctx context.Context, after *string, size *int) (*model.FokontanyList, error) {
	r.FokontanyService.List(0, 100)
	return &model.FokontanyList{}, nil
}

func (r *queryResolver) CountFokontany(ctx context.Context) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Region(ctx context.Context, id string) (*model.Region, error) {
	r.RegionService.ShowRegion(id)
	return &model.Region{}, nil
}

func (r *queryResolver) Regions(ctx context.Context, after *string, size *int) (*model.RegionList, error) {
	r.RegionService.List(0, 100)
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) CountRegions(ctx context.Context) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Search(ctx context.Context, keyword string) (*model.Results, error) {
	r.SearchService.Search(keyword, types.SearchOptions{})
	return &model.Results{}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
