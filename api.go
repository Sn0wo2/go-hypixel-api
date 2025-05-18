package hypixel

import (
	"bytes"
	"io"
	"net/http"
)

type Request struct {
	Method  string
	Header  http.Header
	Path    string
	URL     string // Replace full url in PreRequestHook and Callback!
	Params  Params
	Payload []byte
}

type Response struct {
	Header  http.Header
	Path    string
	URL     string // Replace full url in Callback
	Status  int
	Content []byte
}

// Get Hypixel API HTTP Request
func (c *Client) Get(r Request) (Response, error) {
	if r.Method == "" {
		r.Method = http.MethodGet
	}
	r.URL = c.GetFullPath(r.Path)
	if r.Params == nil {
		r.Params = Params{}
	}
	r.URL = r.Params.String(r.URL)

	if c.GetPreRequestHook() != nil {
		response, err := c.GetPreRequestHook()(r)
		if err == nil {
			return response, nil
		}
	}
	req, err := http.NewRequest(r.Method, r.URL,
		func() io.Reader {
			if r.Payload != nil {
				return bytes.NewReader(r.Payload)
			}
			return nil
		}(),
	)
	if err != nil {
		return Response{}, err
	}
	if r.Header != nil {
		req.Header = r.Header
	}
	if c.GetRate() != nil {
		c.GetRate().WaitIfNeeded()
	}
	rsp, err := c.GetHTTPClient().Do(req)
	if err != nil {
		return Response{}, err
	}
	defer rsp.Body.Close()
	if c.GetRate() != nil {
		_ = c.GetRate().UpdateFromResponse(rsp)
	}
	content, err := io.ReadAll(rsp.Body)
	if err != nil {
		return Response{}, err
	}
	resp := Response{Header: rsp.Header, Path: r.Path, URL: r.URL, Status: rsp.StatusCode, Content: content}
	if c.GetCallback() != nil {
		response, err := c.GetCallback()(r, resp, err)
		if err == nil {
			return response, nil
		}
	}
	return resp, nil
}

// AuthHeader Add api key to header
//
// https://api.hypixel.net/#section/Authentication/ApiKey
func (c *Client) AuthHeader(header ...http.Header) http.Header {
	var h http.Header
	if len(header) == 0 {
		h = http.Header{}
	} else {
		h = header[0]
	}
	h.Set("API-Key", c.apiKey)
	return h
}

// GetPlayerData Data of a specific player, including game stats
// NEED API Key
//
// https://api.hypixel.net/#tag/Player-Data
func (c *Client) GetPlayerData(uuid string) (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "player",
		Params: Params{
			"uuid": uuid,
		},
	})
}

// GetRecentGames The recently played games of a specific player
// NEED API Key
//
// https://api.hypixel.net/#tag/Player-Data/paths/~1v2~1recentgames/get
func (c *Client) GetRecentGames(uuid string) (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "recentgames",
		Params: Params{
			"uuid": uuid,
		},
	})
}

// GetStatus The current online status of a specific player
// NEED API Key
//
// https://api.hypixel.net/#tag/Player-Data/paths/~1v2~1status/get
func (c *Client) GetStatus(uuid string) (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "status",
		Params: Params{
			"uuid": uuid,
		},
	})
}

// GetGuild Retrieve a Guild by a player, id, or name
// NEED API Key
//
// https://api.hypixel.net/#tag/Player-Data/paths/~1v2~1guild/get
func (c *Client) GetGuild(id, player, name string) (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "guild",
		Params: Params{
			"id":     id,
			"player": player,
			"name":   name,
		},
	})
}

// GetGamesInformation Game Information
// Returns information about Hypixel Games. This endpoint is in early development and we are working to add more information when possible
//
// https://api.hypixel.net/#tag/Resources/paths/~1v2~1resources~1games/get
func (c *Client) GetGamesInformation() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Path:   "resources/games",
	})
}

// GetAchievements Achievements
//
// https://api.hypixel.net/#tag/Resources/paths/~1v2~1resources~1achievements/get
func (c *Client) GetAchievements() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Path:   "resources/achievements",
	})
}

// GetChallenges Challenges
//
// https://api.hypixel.net/#tag/Resources/paths/~1v2~1resources~1challenges/get
func (c *Client) GetChallenges() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Path:   "resources/challenges",
	})
}

// GetQuests Quests
//
// https://api.hypixel.net/#tag/Resources/paths/~1v2~1resources~1quests/get
func (c *Client) GetQuests() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Path:   "resources/quests",
	})
}

// GetGuildAchievements Guild Achievements
//
// https://api.hypixel.net/#tag/Resources/paths/~1v2~1resources~1guilds~1achievements/get
func (c *Client) GetGuildAchievements() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Path:   "resources/guilds/achievements",
	})
}

// GetVanityPets Vanity Pets
//
// https://api.hypixel.net/#tag/Resources/paths/~1v2~1resources~1vanity~1pets/get
func (c *Client) GetVanityPets() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Path:   "resources/vanity/pets",
	})
}

// GetVanityCompanions Vanity Companions
//
// https://api.hypixel.net/#tag/Resources/paths/~1v2~1resources~1vanity~1companions/get
func (c *Client) GetVanityCompanions() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Path:   "resources/vanity/companions",
	})
}

// GetSkyBlockCollections Collections
// Information regarding Collections in the SkyBlock game.
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1resources~1skyblock~1collections/get
func (c *Client) GetSkyBlockCollections() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Path:   "resources/skyblock/collections",
	})
}

// GetSkyBlockSkills Skills
// Information regarding skills in the SkyBlock game.
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1resources~1skyblock~1skills/get
func (c *Client) GetSkyBlockSkills() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Path:   "resources/skyblock/skills",
	})
}

// GetSkyBlockItems Items
// Information regarding items in the SkyBlock game.
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1resources~1skyblock~1items/get
func (c *Client) GetSkyBlockItems() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Path:   "resources/skyblock/items",
	})
}

// GetSkyBlockElectionAndMayor Election and Mayor
// Information regarding the current mayor and ongoing election in SkyBlock.
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1resources~1skyblock~1election/get
func (c *Client) GetSkyBlockElectionAndMayor() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Path:   "resources/skyblock/election",
	})
}

// GetSkyBlockCurrentBingoEvent Current Bingo Event
// Information regarding the current bingo event and its goals.
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1resources~1skyblock~1bingo/get
func (c *Client) GetSkyBlockCurrentBingoEvent() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Path:   "resources/skyblock/bingo",
	})
}

// GetSkyBlockNews News
// NEED API Key
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1resources~1skyblock~1news/get
func (c *Client) GetSkyBlockNews() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "skyblock/news",
	})
}

// GetAuctions Request auction(s) by the auction UUID, player UUID, or profile UUID.
// Returns the auctions selected by the provided query. Only one query parameter can be used in a single request, and cannot be filtered by multiple.
// NEED API Key
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1auction/get
func (c *Client) GetAuctions(uuid, player, profile string) (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "skyblock/auction",
		Params: Params{
			"uuid":    uuid,
			"player":  player,
			"profile": profile,
		},
	})
}

// GetActiveAuctions Active auctions
// Returns the currently active auctions sorted by last updated first and paginated.
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1auctions/get
func (c *Client) GetActiveAuctions(page uint) (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Path:   "skyblock/auctions",
		Params: Params{
			"page": page,
		},
	})
}

// GetRecentlyEndedAuctions Recently ended auctions
// SkyBlock auctions which ended in the last 60 seconds.
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1auctions_ended/get
func (c *Client) GetRecentlyEndedAuctions() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Path:   "skyblock/auctions_ended",
	})
}

// GetBazaar Bazaar
// Returns the list of products along with their sell summary, buy summary and quick status.
// Product Description
// The returned product info has 3 main fields:
//
// buy_summary
// sell_summary
// quick_status
// buy_summary and are the current top 30 orders for each transaction type (in-game example: Stock of Stonks).sell_summary
//
// quick_status is a computed summary of the live state of the product (used for advanced mode view in the bazaar):
//
// sellVolume and are the sum of item amounts in all orders.buyVolume
// sellPrice and are the weighted average of the top 2% of orders by volume.buyPrice
// movingWeek is the historic transacted volume from last 7d + live state.
// sellOrders and are the count of active orders. buyOrders
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1bazaar/get
func (c *Client) GetBazaar() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Path:   "skyblock/bazaar",
	})
}

// GetProfileByUUID Profile by UUID
// SkyBlock profile data, such as stats, objectives etc. The data returned can differ depending on the players in-game API settings.
// NEED API Key
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1profile/get
func (c *Client) GetProfileByUUID(profile string) (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "skyblock/profile",
		Params: Params{
			"profile": profile,
		},
	})
}

// GetProfilesByPlayer Profiles by player
// NEED API Key
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1profiles/get
func (c *Client) GetProfilesByPlayer(uuid string) (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "skyblock/profiles",
		Params: Params{
			"uuid": uuid,
		},
	})
}

// GetMuseumData Museum data by profile ID
// SkyBlock museum data for all members of the provided profile. The data returned can differ depending on the players in-game API settings.
// NEED API Key
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1museum/get
func (c *Client) GetMuseumData(profile string) (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "skyblock/museum",
		Params: Params{
			"profile": profile,
		},
	})
}

// GetGardenData Garden data by profile ID
// SkyBlock garden data for the provided profile.
// NEED API Key
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1garden/get
func (c *Client) GetGardenData(profile string) (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "skyblock/garden",
		Params: Params{
			"profile": profile,
		},
	})
}

// GetBingoData Bingo data by player
// Bingo data for participated events of the provided player.
// NEED API Key
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1bingo/get
func (c *Client) GetBingoData(uuid string) (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "skyblock/bingo",
		Params: Params{
			"uuid": uuid,
		},
	})
}

// GetActiveOrUpcomingFireSales Active/Upcoming Fire Sales
// Retrieve the currently active or upcoming Fire Sales for SkyBlock.
//
// https://api.hypixel.net/#tag/SkyBlock/paths/~1v2~1skyblock~1firesales/get
func (c *Client) GetActiveOrUpcomingFireSales() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Path:   "skyblock/firesales",
	})
}

// GetCurrentlyActivePublicHouses currently active public houses.
// This data may be cached for a short period of time.
// NEED API Key
//
// https://api.hypixel.net/#tag/Housing/paths/~1v2~1housing~1active/get
func (c *Client) GetCurrentlyActivePublicHouses() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "housing/active",
	})
}

// GetSpecificHouseInformation Information about a specific house.
// This data may be cached for a short period of time.
// NEED API Key
//
// https://api.hypixel.net/#tag/Housing/paths/~1v2~1housing~1house/get
func (c *Client) GetSpecificHouseInformation(house string) (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "housing/house",
		Params: Params{
			"house": house,
		},
	})
}

// GetSpecificPlayerPublicHouses The public houses for a specific player.//
// This data may be cached for a short period of time.
// NEED API Key
//
// https://api.hypixel.net/#tag/Housing/paths/~1v2~1housing~1houses/get
func (c *Client) GetSpecificPlayerPublicHouses(player string) (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "housing/houses",
		Params: Params{
			"player": player,
		},
	})
}

// GetActiveNetworkBoosters Active Network Boosters
// NEED API Key
//
// https://api.hypixel.net/#tag/Other/paths/~1v2~1boosters/get
func (c *Client) GetActiveNetworkBoosters() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "boosters",
	})
}

// GetCurrentPlayerCounts Current Player Counts
// NEED API Key
//
// https://api.hypixel.net/#tag/Other/paths/~1v2~1counts/get
func (c *Client) GetCurrentPlayerCounts() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "counts",
	})
}

// GetCurrentLeaderboards Current Leaderboards
// NEED API Key
//
// https://api.hypixel.net/#tag/Other/paths/~1v2~1leaderboards/get
func (c *Client) GetCurrentLeaderboards() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "leaderboards",
	})
}

// GetPunishmentStatistics Punishment Statistics
// NEED API Key
//
// https://api.hypixel.net/#tag/Other/paths/~1v2~1punishmentstats/get
func (c *Client) GetPunishmentStatistics() (Response, error) {
	return c.Get(Request{
		Method: http.MethodGet,
		Header: c.AuthHeader(),
		Path:   "punishmentstats",
	})
}
