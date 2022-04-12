package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v43/github"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestGetAuthCodeURL(t *testing.T) {
	t.Run("should return a valid auth code url when state is provided", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)

		oauth2Config := oauth2.Config{
			RedirectURL: "/redirect",
			ClientID:    "testclient",
			Scopes:      []string{"read:user", "user:email", "repo"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "/auth",
				TokenURL: "/token",
			},
		}
		state := "1234abcd"
		mockClientBuilder.EXPECT().GetOAuth2Config().Return(&oauth2Config)

		u := service.GetAuthCodeURL(state)
		assert.Equal(t, "/auth?client_id=testclient&redirect_uri=%2Fredirect&response_type=code&scope=read%3Auser+user%3Aemail+repo&state=1234abcd", u)
	})
}

func TestGetToken(t *testing.T) {
	t.Run("should return github token when authorization code is provided", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		// to get the details of github response structure
		// refer - https://docs.github.com/en/developers/apps/building-oauth-apps/authorizing-oauth-apps#response
		tokenResp := `{
			"access_token": "gho_16C7e42F292c6912E7710c838347Ae178B4a",
			"scope": "repo,gist",
			"token_type": "bearer"
		}`
		router.POST("/token", func(c *gin.Context) {
			c.Data(200, "application/json; charset=utf-8", []byte(tokenResp))
		})
		server := httptest.NewServer(router)
		defer server.Close()
		oauth2Config := oauth2.Config{
			RedirectURL: "/redirect",
			ClientID:    "testclient",
			Scopes:      []string{"read:user", "user:email", "repo"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "/auth",
				TokenURL: server.URL + "/token",
			},
		}
		code := "1234abcd"
		mockClientBuilder.EXPECT().GetOAuth2Config().Return(&oauth2Config)

		token, err := service.GetToken(context.Background(), code)
		tokenJSON, _ := json.Marshal(token)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"access_token":"gho_16C7e42F292c6912E7710c838347Ae178B4a", "expiry":"0001-01-01T00:00:00Z", "token_type":"bearer"}`, string(tokenJSON))
	})

	t.Run("should return error when fetching token from github failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)
		server := httptest.NewServer(nil)
		defer server.Close()

		oauth2Config := oauth2.Config{
			RedirectURL: "/redirect",
			ClientID:    "testclient",
			Scopes:      []string{"read:user", "user:email", "repo"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "/auth",
				TokenURL: server.URL + "/token",
			},
		}
		code := "1234abcd"
		mockClientBuilder.EXPECT().GetOAuth2Config().Return(&oauth2Config)

		_, err := service.GetToken(context.Background(), code)
		assert.Error(t, err)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("should return user when github token is provided", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		// to get the details of github response structure
		// refer - https://docs.github.com/en/rest/reference/users#get-the-authenticated-user
		respJSON := `{
			"id": 1234,
			"login": "johndoe",
			"name": "John Doe",
			"location": "San Francisco",
			"email": "john.doe@example.com",
			"avatar_url": "https://github.com/images/error/octocat_happy.gif"
		}`
		router.GET("/user", func(c *gin.Context) {
			c.Data(200, "application/json; charset=utf-8", []byte(respJSON))
		})
		server := httptest.NewServer(router)
		defer server.Close()
		githubClient := github.NewClient(nil)
		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)

		u, err := service.GetUser(context.Background(), oauth2.Token{})
		userJSON, _ := json.Marshal(u)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"avatar_url":"https://github.com/images/error/octocat_happy.gif", "email":"john.doe@example.com", "id":1234, "location":"San Francisco", "login":"johndoe", "name":"John Doe"}`, string(userJSON))
	})

	t.Run("should return user with primary email when user's email is not publically visible", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		// to get the details of github response structure
		// refer - https://docs.github.com/en/rest/reference/users#get-the-authenticated-user
		userRespJSON := `{
			"id": 1234,
			"login": "johndoe",
			"name": "John Doe",
			"location": "San Francisco",
			"avatar_url": "https://github.com/images/error/octocat_happy.gif"
		}`
		router.GET("/user", func(c *gin.Context) {
			c.Data(200, "application/json; charset=utf-8", []byte(userRespJSON))
		})

		const primailEmail = "john.doe@example.com"
		// to get the details of github response structure
		// refer - https://docs.github.com/en/rest/reference/users#list-email-addresses-for-the-authenticated-user
		emailsRespJSON := fmt.Sprintf(`[
			{
				"email": "%s",
				"primary": true,
				"verified": true,
				"visibility": null
			},
			{
				"email": "alternate.john.doe@example.com",
				"primary": false,
				"verified": false,
				"visibility": "public"
			},
			{
				"email": "another.john.doe@example.com",
				"primary": false,
				"verified": true,
				"visibility": null
			}
		]`, primailEmail)
		router.GET("/user/emails", func(c *gin.Context) {
			c.Data(200, "application/json; charset=utf-8", []byte(emailsRespJSON))
		})
		server := httptest.NewServer(router)
		defer server.Close()
		githubClient := github.NewClient(nil)
		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)

		u, err := service.GetUser(context.Background(), oauth2.Token{})
		userJSON, _ := json.Marshal(u)
		assert.NoError(t, err)
		assert.JSONEq(t, fmt.Sprintf(`{"avatar_url":"https://github.com/images/error/octocat_happy.gif", "email":"%s", "id":1234, "location":"San Francisco", "login":"johndoe", "name":"John Doe"}`, primailEmail), string(userJSON))
	})

	t.Run("should return error when processing user info fails(blank email)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		// to get the details of github response structure
		// refer - https://docs.github.com/en/rest/reference/users#get-the-authenticated-user
		respJSON := `{
			"id": 1234,
			"login": "johndoe",
			"name": "John Doe",
			"location": "San Francisco",
			"avatar_url": "https://github.com/images/error/octocat_happy.gif"
		}`
		router.GET("/user", func(c *gin.Context) {
			c.Data(200, "application/json; charset=utf-8", []byte(respJSON))
		})
		server := httptest.NewServer(router)
		defer server.Close()
		githubClient := github.NewClient(nil)
		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)

		_, err := service.GetUser(context.Background(), oauth2.Token{})
		assert.Error(t, err)
	})

	t.Run("should return error when fetching user fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)
		server := httptest.NewServer(nil)
		defer server.Close()
		githubClient := github.NewClient(nil)
		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)

		_, err := service.GetUser(context.Background(), oauth2.Token{})
		assert.Error(t, err)
	})
}

func TestGetRepos(t *testing.T) {
	t.Run("should return repos when github token is provided", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		server := httptest.NewServer(router)
		defer server.Close()

		// to get the details of github response structure
		// refer - https://docs.github.com/en/rest/reference/repos#list-organization-repositories
		respJSON := `[
			{
				"name": "testrepo",
				"visibility": "private",
				"default_branch": "main"
			},
			{
				"name": "sample-repo",
				"visibility": "public",
				"default_branch": "master"
			}
		]`
		router.GET("/user/repos", func(c *gin.Context) {
			c.Data(200, "application/json; charset=utf-8", []byte(respJSON))
		})
		githubClient := github.NewClient(nil)
		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)

		gitRepos, err := service.GetRepos(context.Background(), oauth2.Token{})
		gitReposJSON, _ := json.Marshal(gitRepos)
		assert.NoError(t, err)
		assert.JSONEq(t, `[{"DefaultBranch":"main", "Name":"testrepo", "Visibility":"private"}, {"DefaultBranch":"master", "Name":"sample-repo", "Visibility":"public"}]`, string(gitReposJSON))
	})

	t.Run("should return error when fetching repos failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)
		server := httptest.NewServer(nil)
		defer server.Close()

		githubClient := github.NewClient(nil)
		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)

		_, err := service.GetRepos(context.Background(), oauth2.Token{})
		assert.Error(t, err)
	})
}

func TestCreateRepo(t *testing.T) {
	t.Run("should create a new repo when request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		server := httptest.NewServer(router)
		defer server.Close()

		// to get the details of github response structure
		// refer - https://docs.github.com/en/rest/reference/repos#create-a-repository-for-the-authenticated-user
		respJSON := `{
			"id": 2250277,
			"node_id": "ZZEwOlJlcG9zaZRvUikXPjk2MjZ5",
			"name": "notes",
			"full_name": "johndoe/notes",
			"owner": {
				"login": "johndoe",
				"id": 1,
				"type": "User"
			},
			"private": true,
			"html_url": "https://github.com/johndoe/notes"
		}`
		router.POST("/user/repos", func(c *gin.Context) {
			c.Data(200, "application/json; charset=utf-8", []byte(respJSON))
		})
		githubClient := github.NewClient(nil)
		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)

		err := service.CreateRepo(context.Background(), oauth2.Token{}, "notes")
		assert.NoError(t, err)
	})

	t.Run("should return error when repo creation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)
		server := httptest.NewServer(nil)
		defer server.Close()

		githubClient := github.NewClient(nil)
		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)

		err := service.CreateRepo(context.Background(), oauth2.Token{}, "notes")
		assert.Error(t, err)
	})
}

func TestSearchFiles(t *testing.T) {
	t.Run("should get search result(files) from git when search request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		server := httptest.NewServer(router)
		defer server.Close()

		// to get the details of github response structure
		// refer - https://docs.github.com/en/rest/reference/search#search-code
		respJSON := `{
			"total_count": 2,
			"incomplete_results": false,
  		"items": [{
				"name": "classes.md",
				"path": "foo/classes.md",
				"sha": "d7212f9dee2dcc18f084d7df8f417b80846ded5a"
			},{
				"name": "birthdays.md",
				"path": "foo/bar/birthdays.md",
				"sha": "c459a67dee2dc4726d2458a32f417699b46da3d9"
			}]
		}`
		router.GET("/search/code", func(c *gin.Context) {
			c.Data(200, "application/json; charset=utf-8", []byte(respJSON))
		})
		router.GET("/repos/testowner/testrepo/contents/foo/classes.md", func(c *gin.Context) {
			c.Data(200, "application/json; charset=utf-8", []byte(`{
				"sha": "d7212f9dee2dcc18f084d7df8f417b80846ded5a",
				"type": "file",
				"size": 5,
				"path": "foo/classes.md",
				"content": "Birthdays"
			}`))
		})
		router.GET("/repos/testowner/testrepo/contents/foo/bar/birthdays.md", func(c *gin.Context) {
			c.Data(200, "application/json; charset=utf-8", []byte(`{
				"sha": "c459a67dee2dc4726d2458a32f417699b46da3d9",
				"type": "file",
				"size": 14,
				"path": "foo/bar/birthdays.md",
				"content": "Birthdays"
			}`))
		})
		fp := GitFileProps{SHA: "", Path: "testpath", Content: "", AuthorName: "", AuthorEmail: "", RepoDetails: GitRepoProps{Repository: "testrepo", DefaultBranch: "", Owner: "testowner"}}
		githubClient := github.NewClient(nil)
		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)

		gitFiles, total, err := service.SearchFiles(context.Background(), oauth2.Token{}, fp, "foo", 1)
		gitFilesJSON, _ := json.Marshal(gitFiles)

		assert.Equal(t, 2, total)
		assert.NoError(t, err)
		assert.JSONEq(t, `[{"Content":"Birthdays", "IsDir":false, "Path":"foo/classes.md", "SHA":"d7212f9dee2dcc18f084d7df8f417b80846ded5a", "Size":5},{"Content":"Birthdays", "IsDir":false, "Path":"foo/bar/birthdays.md", "SHA":"c459a67dee2dc4726d2458a32f417699b46da3d9", "Size":14}]`, string(gitFilesJSON))
	})

	t.Run("should return error when retrieving file info against searched files fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		server := httptest.NewServer(router)
		defer server.Close()

		// to get the details of github response structure
		// refer - https://docs.github.com/en/rest/reference/search#search-code
		respJSON := `{
			"total_count": 2,
			"incomplete_results": false,
  		"items": [{
				"name": "classes.md",
				"path": "foo/classes.md",
				"sha": "d7212f9dee2dcc18f084d7df8f417b80846ded5a"
			},{
				"name": "birthdays.md",
				"path": "foo/bar/birthdays.md",
				"sha": "c459a67dee2dc4726d2458a32f417699b46da3d9"
			}]
		}`
		router.GET("/search/code", func(c *gin.Context) {
			c.Data(200, "application/json; charset=utf-8", []byte(respJSON))
		})
		fp := GitFileProps{SHA: "", Path: "testpath", Content: "", AuthorName: "", AuthorEmail: "", RepoDetails: GitRepoProps{Repository: "testrepo", DefaultBranch: "", Owner: "testowner"}}
		githubClient := github.NewClient(nil)
		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)

		_, _, err := service.SearchFiles(context.Background(), oauth2.Token{}, fp, "foo", 1)
		assert.Error(t, err)
	})

	t.Run("should return error when search request fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)

		server := httptest.NewServer(nil)
		defer server.Close()
		fp := GitFileProps{SHA: "", Path: "testpath", Content: "", AuthorName: "", AuthorEmail: "", RepoDetails: GitRepoProps{Repository: "testrepo", DefaultBranch: "", Owner: "testowner"}}
		githubClient := github.NewClient(nil)
		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)

		_, _, err := service.SearchFiles(context.Background(), oauth2.Token{}, fp, "foo", 1)
		assert.Error(t, err)
	})
}

func TestGetFile(t *testing.T) {
	t.Run("should get the file from git when get request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		server := httptest.NewServer(router)
		defer server.Close()

		// to get the details of github response structure
		// refer - https://docs.github.com/en/rest/reference/repos#get-repository-content
		respJSON := `{
			"sha": "3d21ec53a331a6f037a91c368710b99387d012c1",
			"type": "file",
			"size": 5,
			"path": "testfile.md",
			"content": "Hello"
		}`
		router.GET("/repos/testowner/testrepo/contents/testfile.md", func(c *gin.Context) {
			c.Data(200, "application/json; charset=utf-8", []byte(respJSON))
		})
		fp := GitFileProps{SHA: "", Path: "testfile.md", Content: "", AuthorName: "", AuthorEmail: "", RepoDetails: GitRepoProps{Repository: "testrepo", DefaultBranch: "", Owner: "testowner"}}
		githubClient := github.NewClient(nil)
		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)

		gitFile, err := service.GetFile(context.Background(), oauth2.Token{}, fp)
		gitFileJSON, _ := json.Marshal(gitFile)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"Content":"Hello", "IsDir":false, "Path":"testfile.md", "SHA":"3d21ec53a331a6f037a91c368710b99387d012c1", "Size":5}`, string(gitFileJSON))
	})

	t.Run("should return error when response is of type directory contents", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		server := httptest.NewServer(router)
		defer server.Close()

		// to get the details of github response structure
		// refer - https://docs.github.com/en/rest/reference/repos#get-repository-content
		respJSON := `[{
			"sha": "3d21ec53a331a6f037a91c368710b99387d012c1",
			"type": "dir",
			"size": 0,
			"path": "testpath"
		}]`
		router.GET("/repos/testowner/testrepo/contents/testpath", func(c *gin.Context) {
			c.Data(200, "application/json; charset=utf-8", []byte(respJSON))
		})
		fp := GitFileProps{SHA: "", Path: "testpath", Content: "", AuthorName: "", AuthorEmail: "", RepoDetails: GitRepoProps{Repository: "testrepo", DefaultBranch: "", Owner: "testowner"}}
		githubClient := github.NewClient(nil)

		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)

		_, err := service.GetFile(context.Background(), oauth2.Token{}, fp)
		assert.Error(t, err)
	})

	t.Run("should return error when file is not found at mentioned path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)

		server := httptest.NewServer(nil)
		defer server.Close()

		fp := GitFileProps{SHA: "", Path: "testpath", Content: "", AuthorName: "", AuthorEmail: "", RepoDetails: GitRepoProps{Repository: "testrepo", DefaultBranch: "", Owner: "testowner"}}
		githubClient := github.NewClient(nil)

		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)

		_, err := service.GetFile(context.Background(), oauth2.Token{}, fp)
		assert.Error(t, err)
	})
}

func TestSaveFile(t *testing.T) {
	t.Run("should save a file when save request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		server := httptest.NewServer(router)
		defer server.Close()

		// to get the details of github response structure
		// refer - https://docs.github.com/en/rest/reference/repos#create-or-update-file-contents
		respJSON := `{
			"content": {
				"sha": "3d21ec53a331a6f037a91c368710b99387d012c1",
				"type": "file",
				"size": 5,
				"path": "foo/bar/testfile.md"
			}
		}`
		router.PUT("/repos/testowner/testrepo/contents/foo/bar/testfile.md", func(c *gin.Context) {
			c.Data(200, "application/json; charset=utf-8", []byte(respJSON))
		})
		fp := GitFileProps{SHA: "3d21ec53a331a6f037a91c368710b99387d012c1", Path: "foo/bar/testfile.md", Content: "Hello", AuthorName: "John Doe", AuthorEmail: "john.doe@example.com", RepoDetails: GitRepoProps{Repository: "testrepo", DefaultBranch: "main", Owner: "testowner"}}
		githubClient := github.NewClient(nil)
		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)

		gitFile, err := service.SaveFile(context.Background(), oauth2.Token{}, fp)
		gitFileJSON, _ := json.Marshal(gitFile)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"Content":"", "IsDir":false, "Path":"foo/bar/testfile.md", "SHA":"3d21ec53a331a6f037a91c368710b99387d012c1", "Size":5}`, string(gitFileJSON))
	})

	t.Run("should return error if saving file on github fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)
		server := httptest.NewServer(nil)
		defer server.Close()

		fp := GitFileProps{SHA: "3d21ec53a331a6f037a91c368710b99387d012c1", Path: "foo/bar/testfile.md", Content: "Hello", AuthorName: "John Doe", AuthorEmail: "john.doe@example.com", RepoDetails: GitRepoProps{Repository: "testrepo", DefaultBranch: "main", Owner: "testowner"}}
		githubClient := github.NewClient(nil)
		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)
		_, err := service.SaveFile(context.Background(), oauth2.Token{}, fp)

		assert.Error(t, err)
	})
}

func TestDeleteFile(t *testing.T) {
	t.Run("should delete file when delete request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		server := httptest.NewServer(router)
		defer server.Close()
		router.DELETE("/repos/testowner/testrepo/contents/foo/bar/testfile.md", func(c *gin.Context) {
			c.Data(200, "application/json; charset=utf-8", []byte(""))
		})
		fp := GitFileProps{SHA: "3d21ec53a331a6f037a91c368710b99387d012c1", Path: "foo/bar/testfile.md", Content: "Hello", AuthorName: "John Doe", AuthorEmail: "john.doe@example.com", RepoDetails: GitRepoProps{Repository: "testrepo", DefaultBranch: "main", Owner: "testowner"}}
		githubClient := github.NewClient(nil)
		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)
		err := service.DeleteFile(context.Background(), oauth2.Token{}, fp)

		assert.NoError(t, err)
	})

	t.Run("should return error if deleting file on github fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClientBuilder := NewMockClientBuilder(ctrl)
		service := NewService(mockClientBuilder)
		server := httptest.NewServer(nil)
		defer server.Close()
		fp := GitFileProps{SHA: "3d21ec53a331a6f037a91c368710b99387d012c1", Path: "foo/bar/testfile.md", Content: "Hello", AuthorName: "John Doe", AuthorEmail: "john.doe@example.com", RepoDetails: GitRepoProps{Repository: "testrepo", DefaultBranch: "main", Owner: "testowner"}}
		githubClient := github.NewClient(nil)
		url, _ := url.Parse(server.URL + "/")
		githubClient.BaseURL = url
		mockClientBuilder.EXPECT().Build(gomock.Any(), gomock.Any()).Return(githubClient)
		err := service.DeleteFile(context.Background(), oauth2.Token{}, fp)

		assert.Error(t, err)
	})
}
