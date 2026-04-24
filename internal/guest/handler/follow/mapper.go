package follow

import (
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

func listFollowersRequestToInput(
	req dto.ListFollowersRequest,
	targetCustomerID string,
	requesterCustomerID *string,
) entity.ListFollowersInput {
	return entity.ListFollowersInput{
		TargetCustomerID:    targetCustomerID,
		RequesterCustomerID: requesterCustomerID,
		PageSize:            req.PageSize,
		PageToken:           req.PageToken,
	}
}

func listFollowingRequestToInput(
	req dto.ListFollowingRequest,
	targetCustomerID string,
	requesterUserID *string,
) entity.ListFollowingInput {
	return entity.ListFollowingInput{
		TargetCustomerID:    targetCustomerID,
		RequesterCustomerID: requesterUserID,
		PageSize:            req.PageSize,
		PageToken:           req.PageToken,
	}
}

func listFollowersOutputToResponse(
	out entity.ListFollowersOutput,
) dto.ListCustomersResponse {
	return dto.ListCustomersResponse{
		Customers:     customersToResponse(out.Customers),
		NextPageToken: out.NextPageToken,
	}
}

func listFollowingOutputToResponse(
	out entity.ListFollowingOutput,
) dto.ListCustomersResponse {
	return dto.ListCustomersResponse{
		Customers:     customersToResponse(out.Customers),
		NextPageToken: out.NextPageToken,
	}
}
