package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type videoData struct {
	Data struct {
		XdtShortcodeMedia struct {
			Typename                 string `json:"__typename"`
			IsXDTGraphMediaInterface string `json:"__isXDTGraphMediaInterface"`
			ID                       string `json:"id"`
			Shortcode                string `json:"shortcode"`
			ThumbnailSrc             string `json:"thumbnail_src"`
			Dimensions               struct {
				Height int `json:"height"`
				Width  int `json:"width"`
			} `json:"dimensions"`
			GatingInfo              any `json:"gating_info"`
			FactCheckOverallRating  any `json:"fact_check_overall_rating"`
			FactCheckInformation    any `json:"fact_check_information"`
			SensitivityFrictionInfo any `json:"sensitivity_friction_info"`
			SharingFrictionInfo     struct {
				ShouldHaveSharingFriction bool `json:"should_have_sharing_friction"`
				BloksAppURL               any  `json:"bloks_app_url"`
			} `json:"sharing_friction_info"`
			MediaOverlayInfo any    `json:"media_overlay_info"`
			MediaPreview     string `json:"media_preview"`
			DisplayURL       string `json:"display_url"`
			DisplayResources []struct {
				Src          string `json:"src"`
				ConfigWidth  int    `json:"config_width"`
				ConfigHeight int    `json:"config_height"`
			} `json:"display_resources"`
			AccessibilityCaption any `json:"accessibility_caption"`
			DashInfo             struct {
				IsDashEligible    bool `json:"is_dash_eligible"`
				VideoDashManifest any  `json:"video_dash_manifest"`
				NumberOfQualities int  `json:"number_of_qualities"`
			} `json:"dash_info"`
			HasAudio                  bool    `json:"has_audio"`
			VideoURL                  string  `json:"video_url"`
			VideoViewCount            int     `json:"video_view_count"`
			VideoPlayCount            int     `json:"video_play_count"`
			EncodingStatus            any     `json:"encoding_status"`
			IsPublished               bool    `json:"is_published"`
			ProductType               string  `json:"product_type"`
			Title                     any     `json:"title"`
			VideoDuration             float64 `json:"video_duration"`
			ClipsMusicAttributionInfo struct {
				ArtistName            string `json:"artist_name"`
				SongName              string `json:"song_name"`
				UsesOriginalAudio     bool   `json:"uses_original_audio"`
				ShouldMuteAudio       bool   `json:"should_mute_audio"`
				ShouldMuteAudioReason string `json:"should_mute_audio_reason"`
				AudioID               string `json:"audio_id"`
			} `json:"clips_music_attribution_info"`
			IsVideo               bool   `json:"is_video"`
			TrackingToken         string `json:"tracking_token"`
			UpcomingEvent         any    `json:"upcoming_event"`
			EdgeMediaToTaggedUser struct {
				Edges []any `json:"edges"`
			} `json:"edge_media_to_tagged_user"`
			Owner struct {
				ID                        string `json:"id"`
				Username                  string `json:"username"`
				IsVerified                bool   `json:"is_verified"`
				ProfilePicURL             string `json:"profile_pic_url"`
				BlockedByViewer           bool   `json:"blocked_by_viewer"`
				RestrictedByViewer        any    `json:"restricted_by_viewer"`
				FollowedByViewer          bool   `json:"followed_by_viewer"`
				FullName                  string `json:"full_name"`
				HasBlockedViewer          bool   `json:"has_blocked_viewer"`
				IsEmbedsDisabled          bool   `json:"is_embeds_disabled"`
				IsPrivate                 bool   `json:"is_private"`
				IsUnpublished             bool   `json:"is_unpublished"`
				RequestedByViewer         bool   `json:"requested_by_viewer"`
				PassTieringRecommendation bool   `json:"pass_tiering_recommendation"`
				EdgeOwnerToTimelineMedia  struct {
					Count int `json:"count"`
				} `json:"edge_owner_to_timeline_media"`
				EdgeFollowedBy struct {
					Count int `json:"count"`
				} `json:"edge_followed_by"`
			} `json:"owner"`
			EdgeMediaToCaption struct {
				Edges []any `json:"edges"`
			} `json:"edge_media_to_caption"`
			CanSeeInsightsAsBrand     bool `json:"can_see_insights_as_brand"`
			CaptionIsEdited           bool `json:"caption_is_edited"`
			HasRankedComments         bool `json:"has_ranked_comments"`
			LikeAndViewCountsDisabled bool `json:"like_and_view_counts_disabled"`
			EdgeMediaToParentComment  struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool `json:"has_next_page"`
					EndCursor   any  `json:"end_cursor"`
				} `json:"page_info"`
				Edges []struct {
					Node struct {
						ID              string `json:"id"`
						Text            string `json:"text"`
						CreatedAt       int    `json:"created_at"`
						DidReportAsSpam bool   `json:"did_report_as_spam"`
						Owner           struct {
							ID            string `json:"id"`
							IsVerified    bool   `json:"is_verified"`
							ProfilePicURL string `json:"profile_pic_url"`
							Username      string `json:"username"`
						} `json:"owner"`
						ViewerHasLiked bool `json:"viewer_has_liked"`
						EdgeLikedBy    struct {
							Count int `json:"count"`
						} `json:"edge_liked_by"`
						IsRestrictedPending  bool `json:"is_restricted_pending"`
						EdgeThreadedComments struct {
							Count    int `json:"count"`
							PageInfo struct {
								HasNextPage bool `json:"has_next_page"`
								EndCursor   any  `json:"end_cursor"`
							} `json:"page_info"`
							Edges []any `json:"edges"`
						} `json:"edge_threaded_comments"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_media_to_parent_comment"`
			EdgeMediaToHoistedComment struct {
				Edges []any `json:"edges"`
			} `json:"edge_media_to_hoisted_comment"`
			EdgeMediaPreviewComment struct {
				Count int `json:"count"`
				Edges []struct {
					Node struct {
						ID              string `json:"id"`
						Text            string `json:"text"`
						CreatedAt       int    `json:"created_at"`
						DidReportAsSpam bool   `json:"did_report_as_spam"`
						Owner           struct {
							ID            string `json:"id"`
							IsVerified    bool   `json:"is_verified"`
							ProfilePicURL string `json:"profile_pic_url"`
							Username      string `json:"username"`
						} `json:"owner"`
						ViewerHasLiked bool `json:"viewer_has_liked"`
						EdgeLikedBy    struct {
							Count int `json:"count"`
						} `json:"edge_liked_by"`
						IsRestrictedPending bool `json:"is_restricted_pending"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_media_preview_comment"`
			CommentsDisabled            bool `json:"comments_disabled"`
			CommentingDisabledForViewer bool `json:"commenting_disabled_for_viewer"`
			TakenAtTimestamp            int  `json:"taken_at_timestamp"`
			EdgeMediaPreviewLike        struct {
				Count int   `json:"count"`
				Edges []any `json:"edges"`
			} `json:"edge_media_preview_like"`
			EdgeMediaToSponsorUser struct {
				Edges []any `json:"edges"`
			} `json:"edge_media_to_sponsor_user"`
			IsAffiliate                bool `json:"is_affiliate"`
			IsPaidPartnership          bool `json:"is_paid_partnership"`
			Location                   any  `json:"location"`
			NftAssetInfo               any  `json:"nft_asset_info"`
			ViewerHasLiked             bool `json:"viewer_has_liked"`
			ViewerHasSaved             bool `json:"viewer_has_saved"`
			ViewerHasSavedToCollection bool `json:"viewer_has_saved_to_collection"`
			ViewerInPhotoOfYou         bool `json:"viewer_in_photo_of_you"`
			ViewerCanReshare           bool `json:"viewer_can_reshare"`
			IsAd                       bool `json:"is_ad"`
			EdgeWebMediaToRelatedMedia struct {
				Edges []any `json:"edges"`
			} `json:"edge_web_media_to_related_media"`
			CoauthorProducers   []any `json:"coauthor_producers"`
			PinnedForUsers      []any `json:"pinned_for_users"`
			EdgeRelatedProfiles struct {
				Edges []any `json:"edges"`
			} `json:"edge_related_profiles"`
		} `json:"xdt_shortcode_media"`
	} `json:"data"`
	Extensions struct {
		IsFinal bool `json:"is_final"`
	} `json:"extensions"`
}

func videoDataRequest(id string) (*http.Request, error) {
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	// curl 'https://www.instagram.com/api/graphql' \
	//   -H 'authority: www.instagram.com' \
	//   -H 'accept: */*' \
	//   -H 'accept-language: ru-RU,ru;q=0.9' \
	//   -H 'content-type: application/x-www-form-urlencoded' \
	//   -H 'cookie: csrftoken=KxBv2WNH03JwZdo8gZWE71; mid=ZXIsmQAEAAGUxzE0SXq8Xt9R1sAG; ig_did=FD39D2DC-2696-4BDF-A330-85B3B2C98337; ig_nrcb=1; datr=mSxyZeqF7GH1jjclA2djNOOD; dpr=3' \
	//   -H 'dpr: 3' \
	//   -H 'origin: https://www.instagram.com' \
	//   -H 'referer: https://www.instagram.com/reel/CzMR70XC1C3/' \
	//   -H 'sec-ch-prefers-color-scheme: light' \
	//   -H 'sec-fetch-dest: empty' \
	//   -H 'sec-fetch-mode: cors' \
	//   -H 'sec-fetch-site: same-origin' \
	//   -H 'user-agent: Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1' \
	//   -H 'viewport-width: 390' \
	//   -H 'x-asbd-id: 129477' \
	//   -H 'x-csrftoken: KxBv2WNH03JwZdo8gZWE71' \
	//   -H 'x-fb-friendly-name: PolarisPostActionLoadPostQueryQuery' \
	//   -H 'x-fb-lsd: AVp8-FMus6c' \
	//   -H 'x-ig-app-id: 1217981644879628' \
	//   --data-raw 'av=0&__d=www&__user=0&__a=1&__req=3&__hs=19698.HYP%3Ainstagram_web_pkg.2.1..0.0&dpr=2&__ccg=UNKNOWN&__rev=1010272456&__s=5o02yo%3A9d9qeh%3Abvvaqt&__hsi=7309954327711794454&__dyn=7xeUjG1mxu1syUbFp60DU98nwgU29zEdEc8co2qwJw5ux609vCwjE1xoswIwuo2awlU-cw5Mx62G3i1ywOwv89k2C1Fwc60AEC7U2czXwae4UaEW2G1NwwwNwKwHw8Xxm16wUwtEvw4JwJwSyES1Twoob82cwMwrUdUbGwmk1xwmo6O1FwlE6OFA6fxy4Ujw&__csr=goOQyhq6eExVTucCpb-_nKjCjKEzBVaF4yQqmKdhVVpFWAiADzpXK9BgyRBxCh2JqiDKi-swyhpF9lDrgyVbyi38OmaLUNqyqBhUG8CBxe2y2268nw04Kig2TJw1iOEggjw1X6062E4i3FG44U_c0O9E0Sq2Kh0fnc4U5Wi260vQMy0u20uK00GdE&__comet_req=7&lsd=AVp8-FMus6c&jazoest=2896&__spin_r=1010272456&__spin_b=trunk&__spin_t=1701981371&fb_api_caller_class=RelayModern&fb_api_req_friendly_name=PolarisPostActionLoadPostQueryQuery&variables=%7B%22shortcode%22%3A%22CzMR70XC1C3%22%2C%22fetch_comment_count%22%3A40%2C%22fetch_related_profile_media_count%22%3A3%2C%22parent_comment_count%22%3A24%2C%22child_comment_count%22%3A3%2C%22fetch_like_count%22%3A10%2C%22fetch_tagged_user_count%22%3Anull%2C%22fetch_preview_comment_count%22%3A2%2C%22has_threaded_comments%22%3Atrue%2C%22hoisted_comment_id%22%3Anull%2C%22hoisted_reply_id%22%3Anull%7D&server_timestamps=true&doc_id=10015901848480474' \
	//

	params := url.Values{}
	params.Add("av", `0`)
	params.Add("__d", `www`)
	params.Add("__user", `0`)
	params.Add("__a", `1`)
	params.Add("__req", `3`)
	params.Add("__hs", `19698.HYP:instagram_web_pkg.2.1..0.0`)
	params.Add("dpr", `2`)
	params.Add("__ccg", `UNKNOWN`)
	params.Add("__rev", `1010272456`)
	params.Add("__s", `5o02yo:9d9qeh:bvvaqt`)
	params.Add("__hsi", `7309954327711794454`)
	params.Add(
		"__dyn",
		`7xeUjG1mxu1syUbFp60DU98nwgU29zEdEc8co2qwJw5ux609vCwjE1xoswIwuo2awlU-cw5Mx62G3i1ywOwv89k2C1Fwc60AEC7U2czXwae4UaEW2G1NwwwNwKwHw8Xxm16wUwtEvw4JwJwSyES1Twoob82cwMwrUdUbGwmk1xwmo6O1FwlE6OFA6fxy4Ujw`,
	)
	params.Add(
		"__csr",
		`goOQyhq6eExVTucCpb-_nKjCjKEzBVaF4yQqmKdhVVpFWAiADzpXK9BgyRBxCh2JqiDKi-swyhpF9lDrgyVbyi38OmaLUNqyqBhUG8CBxe2y2268nw04Kig2TJw1iOEggjw1X6062E4i3FG44U_c0O9E0Sq2Kh0fnc4U5Wi260vQMy0u20uK00GdE`,
	)
	params.Add("__comet_req", `7`)
	params.Add("lsd", `AVp8-FMus6c`)
	params.Add("jazoest", `2896`)
	params.Add("__spin_r", `1010272456`)
	params.Add("__spin_b", `trunk`)
	params.Add("__spin_t", `1701981371`)
	params.Add("fb_api_caller_class", `RelayModern`)
	params.Add("fb_api_req_friendly_name", `PolarisPostActionLoadPostQueryQuery`)
	params.Add(
		"variables",
		fmt.Sprintf(
			`{"shortcode":"%s","fetch_comment_count":40,"fetch_related_profile_media_count":3,"parent_comment_count":24,"child_comment_count":3,"fetch_like_count":10,"fetch_tagged_user_count":null,"fetch_preview_comment_count":2,"has_threaded_comments":true,"hoisted_comment_id":null,"hoisted_reply_id":null}`,
			id,
		),
	)
	params.Add("server_timestamps", `true`)
	params.Add("doc_id", `10015901848480474`)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://www.instagram.com/api/graphql", body)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authority", "www.instagram.com")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(
		"Cookie",
		"csrftoken=KxBv2WNH03JwZdo8gZWE71; mid=ZXIsmQAEAAGUxzE0SXq8Xt9R1sAG; ig_did=FD39D2DC-2696-4BDF-A330-85B3B2C98337; ig_nrcb=1; datr=mSxyZeqF7GH1jjclA2djNOOD; dpr=3",
	)
	req.Header.Set("Dpr", "3")
	req.Header.Set("Origin", "https://www.instagram.com")
	req.Header.Set("Referer", fmt.Sprintf("https://www.instagram.com/reel/%s/", id))
	req.Header.Set("Sec-Ch-Prefers-Color-Scheme", "light")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
	)
	req.Header.Set("Viewport-Width", "390")
	req.Header.Set("X-Asbd-Id", "129477")
	req.Header.Set("X-Csrftoken", "KxBv2WNH03JwZdo8gZWE71")
	req.Header.Set("X-Fb-Friendly-Name", "PolarisPostActionLoadPostQueryQuery")
	req.Header.Set("X-Fb-Lsd", "AVp8-FMus6c")
	req.Header.Set("X-Ig-App-Id", "1217981644879628")

	return req, nil
}

func idFromURL(urlStr string) (string, error) {
	_url, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("parse url: %w", err)
	}
	path := strings.Trim(_url.Path, "/")
	parts := strings.Split(path, "/")

	return parts[len(parts)-1], nil
}

func videoDownloadUrl(_url string) (string, error) {
	id, err := idFromURL(_url)
	if err != nil {
		return "", fmt.Errorf("id from url: %w", err)
	}
	req, err := videoDataRequest(id)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code: %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}
	vd := videoData{}
	err = json.Unmarshal(b, &vd)
	if err != nil {
		return "", fmt.Errorf("unmarshal: %w", err)
	}
	videoUrl := vd.Data.XdtShortcodeMedia.VideoURL
	if videoUrl == "" {
		return "", fmt.Errorf("video url is empty")
	}
	return videoUrl, nil
}

func downloadVideo(_url string) error {
	req, err := http.NewRequest(http.MethodGet, _url, nil)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	f, err := os.Create("video.mp4")
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	return nil
}

func main() {
	instagramUrl := ""
	flag.StringVar(&instagramUrl, "url", "", "url")
	flag.Parse()

	if instagramUrl == "" {
		fmt.Println("url is empty")
		os.Exit(1)
	}
	vdu, err := videoDownloadUrl(instagramUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = downloadVideo(vdu)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("ok")
}
