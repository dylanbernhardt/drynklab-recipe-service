package recipe

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
	"time"

	pb "github.com/dylanbernhardt/drynklab-recipe-service/proto/recipe"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	pb.UnimplementedRecipeServiceServer
	mu      sync.RWMutex
	recipes map[int32]*pb.Recipe
	nextID  int32
}

func NewService() *Service {
	return &Service{
		recipes: make(map[int32]*pb.Recipe),
		nextID:  1,
	}
}

func (s *Service) GetRecipe(ctx context.Context, req *pb.GetRecipeRequest) (*pb.Recipe, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "request was canceled")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	recipe, exists := s.recipes[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "recipe with ID %d not found", req.Id)
	}

	return recipe, nil
}

func (s *Service) ListRecipes(ctx context.Context, req *pb.ListRecipesRequest) (*pb.ListRecipesResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "request was canceled")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize
	if start >= int32(len(s.recipes)) {
		return &pb.ListRecipesResponse{Recipes: []*pb.Recipe{}, TotalCount: int32(len(s.recipes))}, nil
	}
	if end > int32(len(s.recipes)) {
		end = int32(len(s.recipes))
	}

	result := make([]*pb.Recipe, 0, end-start)
	for i := start; i < end; i++ {
		result = append(result, s.recipes[i+1])
	}

	return &pb.ListRecipesResponse{
		Recipes:    result,
		TotalCount: int32(len(s.recipes)),
	}, nil
}

func (s *Service) CreateRecipe(ctx context.Context, req *pb.CreateRecipeRequest) (*pb.Recipe, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "request was canceled")
	}

	if err := validateRecipe(req.Name, req.Instructions, req.Ingredients); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := timestamppb.New(time.Now())
	recipe := &pb.Recipe{
		Id:           s.nextID,
		Name:         req.Name,
		Instructions: req.Instructions,
		Ingredients:  req.Ingredients,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	s.recipes[s.nextID] = recipe
	s.nextID++

	return recipe, nil
}

func (s *Service) UpdateRecipe(ctx context.Context, req *pb.UpdateRecipeRequest) (*pb.Recipe, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "request was canceled")
	}

	if err := validateRecipe(req.Recipe.Name, req.Recipe.Instructions, req.Recipe.Ingredients); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	recipe, exists := s.recipes[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "recipe with ID %d not found", req.Id)
	}

	recipe.Name = req.Recipe.Name
	recipe.Instructions = req.Recipe.Instructions
	recipe.Ingredients = req.Recipe.Ingredients
	recipe.UpdatedAt = timestamppb.New(time.Now())

	return recipe, nil
}

func (s *Service) DeleteRecipe(ctx context.Context, req *pb.DeleteRecipeRequest) (*pb.DeleteRecipeResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "request was canceled")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.recipes[req.Id]; !exists {
		return nil, status.Errorf(codes.NotFound, "recipe with ID %d not found", req.Id)
	}

	delete(s.recipes, req.Id)

	return &pb.DeleteRecipeResponse{Success: true}, nil
}

func validateRecipe(name, instructions string, ingredients []*pb.Ingredient) error {
	if name == "" {
		return fmt.Errorf("recipe name cannot be empty")
	}
	if instructions == "" {
		return fmt.Errorf("recipe instructions cannot be empty")
	}
	if len(ingredients) == 0 {
		return fmt.Errorf("recipe must have at least one ingredient")
	}
	for i, ing := range ingredients {
		if ing.Name == "" {
			return fmt.Errorf("ingredient %d name cannot be empty", i+1)
		}
		if ing.Quantity <= 0 {
			return fmt.Errorf("ingredient %d quantity must be positive", i+1)
		}
		if ing.Unit == "" {
			return fmt.Errorf("ingredient %d unit cannot be empty", i+1)
		}
	}
	return nil
}
