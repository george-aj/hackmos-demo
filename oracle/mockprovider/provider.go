package mockprovider

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	distTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"strings"
)

type Provider interface {
	GetDelegatorDelegations(query stakingTypes.QueryDelegatorDelegationsRequest) (stakingTypes.QueryDelegatorDelegationsResponse, error)
	GetDelegatorRewards(query distTypes.QueryDelegationTotalRewardsRequest) (distTypes.QueryDelegationTotalRewardsResponse, error)
	GetWasmQuery() (DataResponse, error)
}

var _ Provider = MockProvider{}

type MockProvider struct{}

func NewMockProvider() MockProvider {
	return MockProvider{}
}

func grpcConn() (grpc.ClientConn, error) {
	// TODO allow each validator to set their own
	grpcUrl := "juno-grpc.polkachu.com:12690"

	var transport grpc.DialOption
	if strings.Contains(grpcUrl, "https://") {
		config := &tls.Config{InsecureSkipVerify: true}
		creds := credentials.NewTLS(config)
		transport = grpc.WithTransportCredentials(creds)
	} else {
		transport = grpc.WithInsecure()
	}

	grpcConn, err := grpc.Dial(grpcUrl, transport)
	if err != nil {
		return grpc.ClientConn{}, err
	}

	return *grpcConn, nil
}

func (p MockProvider) GetDelegatorRewards(query distTypes.QueryDelegationTotalRewardsRequest) (distTypes.QueryDelegationTotalRewardsResponse, error) {
	fmt.Println("Query rewards")
	grpcConn, err := grpcConn()
	if err != nil {
		return distTypes.QueryDelegationTotalRewardsResponse{}, err
	}

	defer grpcConn.Close()

	distQueryClient := distTypes.NewQueryClient(&grpcConn)
	distQueryReq, err := distQueryClient.DelegationTotalRewards(
		context.Background(),
		&query,
	)
	if err != nil {
		return distTypes.QueryDelegationTotalRewardsResponse{}, err
	}

	return *distQueryReq, err
}

func (p MockProvider) GetDelegatorDelegations(query stakingTypes.QueryDelegatorDelegationsRequest) (stakingTypes.QueryDelegatorDelegationsResponse, error) {
	fmt.Println("Query delegations")
	grpcConn, err := grpcConn()
	if err != nil {
		return stakingTypes.QueryDelegatorDelegationsResponse{}, err
	}

	defer grpcConn.Close()

	stakingQueryClient := stakingTypes.NewQueryClient(&grpcConn)
	stakingQueryReq, err := stakingQueryClient.DelegatorDelegations(
		context.Background(),
		&query,
	)
	if err != nil {
		return stakingTypes.QueryDelegatorDelegationsResponse{}, err
	}

	return *stakingQueryReq, err
}

type DataResponse struct {
	Data struct {
		Tokens []string `json:"tokens"`
	} `json:"data"`
}

func (p MockProvider) GetWasmQuery() (DataResponse, error) {
	response, err := http.Get("https://juno-api.lavenderfive.com/cosmwasm/wasm/v1/contract/juno1anh4pf98fe8uh64uuhaasqdmg89qe6kk5xsklxuvtjmu6rhpg53sj9uejj/smart/eyJhbGxfdG9rZW5zIjp7fX0K")

	if err != nil {
		return DataResponse{}, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return DataResponse{}, err
	}

	data := DataResponse{}

	err = json.Unmarshal(responseData, &data)
	if err != nil {
		return DataResponse{}, err
	}

	return data, nil
}
