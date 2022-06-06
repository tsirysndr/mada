package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/tsirysndr/mada/graph/generated"
	"github.com/tsirysndr/mada/graph/model"
	"github.com/tsirysndr/mada/types"
)

func (r *mutationResolver) Empty(ctx context.Context) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Commune(ctx context.Context, id string) (*model.Commune, error) {
	result, err := r.CommuneService.ShowCommune(id)
	if err != nil {
		return nil, err
	}
	var commune model.Commune
	err = copier.Copy(&commune, &result)
	return &commune, err
}

func (r *queryResolver) Communes(ctx context.Context, skip *int, size *int) (*model.CommuneList, error) {
	result, err := r.CommuneService.List(ParseSkipAndSize(skip, size))
	if err != nil {
		return nil, err
	}
	communes := []*model.Commune{}
	for _, v := range result.Hits {
		name := v.Fields["name"].(string)
		district := v.Fields["district"].(string)
		region := v.Fields["region"].(string)
		country := v.Fields["country"].(string)
		communes = append(communes, &model.Commune{
			ID:       &v.ID,
			Name:     &name,
			District: &district,
			Region:   &region,
			Country:  &country,
		})
	}
	return &model.CommuneList{Data: communes}, nil
}

func (r *queryResolver) CountCommunes(ctx context.Context) (int, error) {
	return r.CommuneService.Count()
}

func (r *queryResolver) District(ctx context.Context, id string) (*model.District, error) {
	result, err := r.DistrictService.ShowDistrict(id)
	if err != nil {
		return nil, err
	}
	var district model.District
	err = copier.Copy(&district, &result)
	return &district, err
}

func (r *queryResolver) Districts(ctx context.Context, skip *int, size *int) (*model.DistrictList, error) {
	result, err := r.DistrictService.List(ParseSkipAndSize(skip, size))
	if err != nil {
		return nil, err
	}
	districts := []*model.District{}
	for _, v := range result.Hits {
		name := v.Fields["name"].(string)
		region := v.Fields["region"].(string)
		country := v.Fields["country"].(string)
		districts = append(districts, &model.District{
			ID:      &v.ID,
			Name:    &name,
			Region:  &region,
			Country: &country,
		})
	}
	return &model.DistrictList{Data: districts}, nil
}

func (r *queryResolver) CountDistricts(ctx context.Context) (int, error) {
	return r.DistrictService.Count()
}

func (r *queryResolver) Fokontany(ctx context.Context, id string) (*model.Fokontany, error) {
	result, err := r.FokontanyService.ShowFokontany(id)
	if err != nil {
		return nil, err
	}
	var fokontany model.Fokontany
	err = copier.Copy(&fokontany, &result)
	return &fokontany, err
}

func (r *queryResolver) AllFokontany(ctx context.Context, skip *int, size *int) (*model.FokontanyList, error) {
	result, err := r.FokontanyService.List(ParseSkipAndSize(skip, size))
	if err != nil {
		return nil, err
	}
	list := []*model.Fokontany{}
	for _, v := range result.Hits {
		name := v.Fields["name"].(string)
		commune := v.Fields["commune"].(string)
		district := v.Fields["district"].(string)
		region := v.Fields["region"].(string)
		country := v.Fields["country"].(string)
		list = append(list, &model.Fokontany{
			ID:       &v.ID,
			Name:     &name,
			Commune:  &commune,
			District: &district,
			Region:   &region,
			Country:  &country,
		})
	}
	return &model.FokontanyList{Data: list}, nil
}

func (r *queryResolver) CountFokontany(ctx context.Context) (int, error) {
	return r.FokontanyService.Count()
}

func (r *queryResolver) Region(ctx context.Context, id string) (*model.Region, error) {
	result, err := r.RegionService.ShowRegion(id)
	if err != nil {
		return nil, err
	}
	var region model.Region
	err = copier.Copy(&region, &result)
	return &region, err
}

func (r *queryResolver) Regions(ctx context.Context, skip *int, size *int) (*model.RegionList, error) {
	result, err := r.RegionService.List(ParseSkipAndSize(skip, size))
	if err != nil {
		return nil, err
	}
	regions := []*model.Region{}
	for _, v := range result.Hits {
		name := v.Fields["name"].(string)
		country := v.Fields["country"].(string)
		regions = append(regions, &model.Region{
			ID:      &v.ID,
			Name:    &name,
			Country: &country,
		})
	}
	return &model.RegionList{Data: regions}, nil
}

func (r *queryResolver) CountRegions(ctx context.Context) (int, error) {
	return r.RegionService.Count()
}

func (r *queryResolver) Search(ctx context.Context, keyword string) (*model.Results, error) {
	result, err := r.SearchService.Search(keyword, types.SearchOptions{})
	if err != nil {
		return nil, err
	}
	var (
		commune   model.Commune
		district  model.District
		fokontany model.Fokontany
		region    model.Region
		hits      []*model.Hit
	)
	if result.Result != nil {
		for _, v := range result.Result.Hits {
			var hit model.Hit
			copier.Copy(&hit, &v)
			fokontany := v.Fields["fokontany"].(string)
			commune := v.Fields["commune"].(string)
			district := v.Fields["district"].(string)
			region := v.Fields["region"].(string)
			country := v.Fields["country"].(string)
			_type := v.Fields["type"].(string)
			hit.Fields = &model.Fields{
				Fokontany: &fokontany,
				Commune:   &commune,
				District:  &district,
				Region:    &region,
				Country:   &country,
				Type:      &_type,
			}
			hits = append(hits, &hit)
		}
	}
	copier.Copy(&commune, &result.Commune)
	copier.Copy(&district, &result.District)
	copier.Copy(&fokontany, &result.Fokontany)
	copier.Copy(&region, &result.Region)
	return &model.Results{
		Commune:   &commune,
		District:  &district,
		Fokontany: &fokontany,
		Region:    &region,
		Hits:      hits,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
