package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSaveUserGraphQLSuccess1(t *testing.T) {
	query := `mutation {SaveUser(input: {email: \"test@test.com\", password: \"1234567890\"}) {user_id}}`

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultSaveUserSuccess); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(resultSaveUserSuccess.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", resultSaveUserSuccess.Errors)
	}

	expectedUserID := 1
	if resultSaveUserSuccess.Data.SaveUser.UserID != expectedUserID {
		t.Errorf("unexpected user_id: got %v want %v", resultSaveUserSuccess.Data.SaveUser.UserID, expectedUserID)
	} else {
		fmt.Printf("User registered with ID: %v\n", resultSaveUserSuccess.Data.SaveUser.UserID)
	}
}

func TestSaveUserGraphQLSuccess2(t *testing.T) {
	query := `mutation {SaveUser(input: {email: \"test2@test.com\", password: \"1234567890\"}) {user_id}}`

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultSaveUserSuccess); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(resultSaveUserSuccess.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", resultSaveUserSuccess.Errors)
	}

	expectedUserID := 2
	if resultSaveUserSuccess.Data.SaveUser.UserID != expectedUserID {
		t.Errorf("unexpected user_id: got %v want %v", resultSaveUserSuccess.Data.SaveUser.UserID, expectedUserID)
	} else {
		fmt.Printf("User registered with ID: %v\n", resultSaveUserSuccess.Data.SaveUser.UserID)
	}
}

func TestSaveUserGraphQLFail(t *testing.T) {
	query := `mutation {SaveUser(input: {email: \"test@test.com\", password: \"1234567890\"}) {user_id}}`

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("response Body:", resp.Body)
	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultSaveUserFail); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	expectedUserID := 0
	if resultSaveUserFail.Data.SaveUser.UserID != expectedUserID {
		t.Errorf("unexpected user_id: got %v want %v", resultSaveUserFail.Data.SaveUser.UserID, expectedUserID)
	} else {
		fmt.Printf("User already exists\n")
	}
}

func TestLoginGraphQLSuccess(t *testing.T) {
	query := `mutation {Login(input: {email: \"test@test.com\", password: \"1234567890\"}) {token} }`

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultLoginSuccess); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(resultLoginSuccess.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", resultLoginSuccess.Errors)
	}

	if resultLoginSuccess.Data.Login.Token == "" {
		t.Errorf("unexpected token: got empty string")
	} else {
		fmt.Printf("User logged in with token: %v\n", resultLoginSuccess.Data.Login.Token)
	}
}

func TestLoginGraphQLFail(t *testing.T) {
	query := `mutation {Login(input: {email: \"doesntexists@test.com\", password: \"1234567890\"}) {token} }`

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultLoginFail); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if resultLoginFail.Data.Login.Token != "" {
		t.Errorf("unexpected token:  %v\n", resultLoginFail.Data.Login.Token)
	} else {
		fmt.Printf("User not found\n")
	}
}

func TestSavePostGraphQLSuccess1(t *testing.T) {
	query := fmt.Sprintf(`mutation {SavePost(input: {token: \"%s\", title: \"Test1\", content: \"hello, world\", comments: true}) {post_id created_at}}`, resultLoginSuccess.Data.Login.Token)

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultSavePostSuccess); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(resultSavePostSuccess.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", resultSavePostSuccess.Errors)
	}

	if resultSavePostSuccess.Data.SavePostResponse.PostID == 0 {
		t.Errorf("unexpected post_id: got %v", resultSavePostSuccess.Data.SavePostResponse.PostID)
	} else {
		fmt.Printf("Post saved with ID: %v\n, created: %v\n",
			resultSavePostSuccess.Data.SavePostResponse.PostID, resultSavePostSuccess.Data.SavePostResponse.CreatedAt)
	}
}

func TestSavePostGraphQLSuccess2(t *testing.T) {
	query := fmt.Sprintf(`mutation {SavePost(input: {token: \"%s\", title: \"Test2\", content: \"hello, world\", comments: false}) {post_id created_at}}`, resultLoginSuccess.Data.Login.Token)

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultSavePostSuccess); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(resultSavePostSuccess.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", resultSavePostSuccess.Errors)
	}

	if resultSavePostSuccess.Data.SavePostResponse.PostID == 0 {
		t.Errorf("unexpected post_id: got %v", resultSavePostSuccess.Data.SavePostResponse.PostID)
	} else {
		fmt.Printf("Post saved with ID: %v\n, created: %v\n",
			resultSavePostSuccess.Data.SavePostResponse.PostID, resultSavePostSuccess.Data.SavePostResponse.CreatedAt)
	}
}

func TestSavePostGraphQLFail(t *testing.T) {
	query := fmt.Sprintf(`mutation {SavePost(input: {token: \"%s\", title: \"Test2\", content: \"hello, world\", comments: false}) {post_id created_at}}`, resultLoginSuccess.Data.Login.Token)

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultSavePostFail); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if resultSavePostFail.Data.SavePostResponse.PostID != 0 {
		t.Errorf("unexpected post_id: got %v", resultSavePostFail.Data.SavePostResponse.PostID)
	} else {
		fmt.Printf("Post already exists\n")
	}
}

func TestProvidePostGraphQLSuccess1(t *testing.T) {
	query := `mutation {ProvidePost(input: {post_id: 1}) {id user_id title content created_at comments} }`

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultProvidePostSuccess); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(resultProvidePostSuccess.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", resultProvidePostSuccess.Errors)
	}

	if resultProvidePostSuccess.Data.ProvidePostResponse.ID == 0 {
		t.Errorf("unexpected post_id: got %v", resultProvidePostSuccess.Data.ProvidePostResponse.ID)
	} else {
		fmt.Printf("Post received: %v, title: %v, content: %v\n",
			resultProvidePostSuccess.Data.ProvidePostResponse.ID,
			resultProvidePostSuccess.Data.ProvidePostResponse.Title,
			resultProvidePostSuccess.Data.ProvidePostResponse.Content)
	}
}

func TestProvidePostGraphQLSuccess2(t *testing.T) {
	query := `mutation {ProvidePost(input: {post_id: 2}) {id user_id title content created_at comments} }`

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultProvidePostSuccess); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(resultProvidePostSuccess.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", resultProvidePostSuccess.Errors)
	}

	if resultProvidePostSuccess.Data.ProvidePostResponse.ID == 0 {
		t.Errorf("unexpected post_id: got %v", resultProvidePostSuccess.Data.ProvidePostResponse.ID)
	} else {
		fmt.Printf("Post received: %v, title: %v, content: %v\n",
			resultProvidePostSuccess.Data.ProvidePostResponse.ID,
			resultProvidePostSuccess.Data.ProvidePostResponse.Title,
			resultProvidePostSuccess.Data.ProvidePostResponse.Content)
	}
}

func TestProvidePostGraphQLFail(t *testing.T) {
	query := `mutation {ProvidePost(input: {post_id: 3}) {id user_id title content created_at comments} }`

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultProvidePostFail); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if resultProvidePostFail.Data.ProvidePostResponse.ID != 0 {
		t.Errorf("unexpected post_id: got %v", resultProvidePostFail.Data.ProvidePostResponse.ID)
	} else {
		fmt.Printf("Post not found\n")
	}
}

func TestProvideAllPostsGraphQL1(t *testing.T) {
	query := `mutation { ProvideAllPosts(input: {page: 1}) {posts {id user_id title content created_at comments} } }`

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	fmt.Println()
	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultProvideAllPosts); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(resultProvideAllPosts.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", resultProvideAllPosts.Errors)
	}

	if len(resultProvideAllPosts.Data.ProvideAllPostsResponse.Posts) == 0 {
		fmt.Printf("No posts found\n")
	} else {
		fmt.Printf("Posts received: %v\n",
			resultProvideAllPosts.Data.ProvideAllPostsResponse.Posts)
	}
}

func TestProvideAllPostsGraphQL2(t *testing.T) {
	query := `mutation { ProvideAllPosts(input: {page: 1000}) {posts {id user_id title content created_at comments} } }`

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	fmt.Println()
	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultProvideAllPosts); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(resultProvideAllPosts.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", resultProvideAllPosts.Errors)
	}

	if len(resultProvideAllPosts.Data.ProvideAllPostsResponse.Posts) == 0 {
		fmt.Printf("No posts found\n")
	} else {
		fmt.Printf("Posts received: %v\n",
			resultProvideAllPosts.Data.ProvideAllPostsResponse.Posts)
	}
}

func TestSaveCommentGraphQLSuccess1(t *testing.T) {
	query := fmt.Sprintf(`mutation {SaveComment(input: {token:\"%s\", post_id: 1, content: \"hello\"}) {id created_at} }`, resultLoginSuccess.Data.Login.Token)

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	fmt.Println()
	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultSaveCommentSuccess); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(resultSaveCommentSuccess.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", resultSaveCommentSuccess.Errors)
	}

	expectedCommentID := 1
	if resultSaveCommentSuccess.Data.SaveCommentResponse.ID != expectedCommentID {
		t.Errorf("unexpected comment_id: got %v want %v", resultSaveCommentSuccess.Data.SaveCommentResponse.ID, expectedCommentID)
	} else {
		fmt.Printf("Comment saved with ID: %v\n", resultSaveCommentSuccess.Data.SaveCommentResponse.ID)
	}
}

func TestSaveCommentGraphQLSuccess2(t *testing.T) {
	query := fmt.Sprintf(`mutation {SaveComment(input: {token:\"%s\", post_id: 1, content: \"hello2\"}) {id created_at} }`, resultLoginSuccess.Data.Login.Token)

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	fmt.Println()
	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultSaveCommentSuccess); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(resultSaveCommentSuccess.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", resultSaveCommentSuccess.Errors)
	}

	expectedCommentID := 2
	if resultSaveCommentSuccess.Data.SaveCommentResponse.ID != expectedCommentID {
		t.Errorf("unexpected comment_id: got %v want %v", resultSaveCommentSuccess.Data.SaveCommentResponse.ID, expectedCommentID)
	} else {
		fmt.Printf("Comment saved with ID: %v\n", resultSaveCommentSuccess.Data.SaveCommentResponse.ID)
	}
}

func TestSaveCommentGraphQLFail1(t *testing.T) {
	query := fmt.Sprintf(`mutation {SaveComment(input: {token:\"%s\", post_id: 2, content: \"hello2\"}) {id created_at} }`, resultLoginSuccess.Data.Login.Token)

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	fmt.Println()
	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultSaveCommentFail); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if resultSaveCommentFail.Data.SaveCommentResponse.ID != 0 {
		t.Errorf("unexpected comment_id: got %v", resultSaveCommentFail.Data.SaveCommentResponse.ID)
	} else {
		fmt.Printf("Comments are not allowed\n")
	}
}

func TestSaveCommentGraphQLFail2(t *testing.T) {
	query := fmt.Sprintf(`mutation {SaveComment(input: {token:\"%s\", post_id: 2, content: \"hello2\"}) {id created_at} }`, resultLoginFail.Data.Login.Token)

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	fmt.Println()
	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultSaveCommentFail); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if resultSaveCommentFail.Data.SaveCommentResponse.ID != 0 {
		t.Errorf("unexpected comment_id: got %v", resultSaveCommentFail.Data.SaveCommentResponse.ID)
	} else {
		fmt.Printf("Access denied\n")
	}
}

func TestSaveCommentToCommentGraphQLSuccess1(t *testing.T) {
	query := fmt.Sprintf(`mutation {SaveCommentToComment(input: {token:\"%s\", post_id: 1, parent_id: 1, content: \"hello comment to comment\"}) {id created_at} }`, resultLoginSuccess.Data.Login.Token)

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	fmt.Println()
	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultSaveCommentToCommentSuccess); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(resultSaveCommentToCommentSuccess.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", resultSaveCommentToCommentSuccess.Errors)
	}

	expectedCommentID := 3
	if resultSaveCommentToCommentSuccess.Data.SaveCommentResponse.ID != expectedCommentID {
		t.Errorf("unexpected comment_id: got %v want %v", resultSaveCommentToCommentSuccess.Data.SaveCommentResponse.ID, expectedCommentID)
	} else {
		fmt.Printf("Comment to comment saved with ID: %v\n", resultSaveCommentToCommentSuccess.Data.SaveCommentResponse.ID)
	}
}

func TestSaveCommentToCommentGraphQLSuccess2(t *testing.T) {
	query := fmt.Sprintf(`mutation {SaveCommentToComment(input: {token:\"%s\", post_id: 1, parent_id: 1, content: \"hello comment to comment 2\"}) {id created_at} }`, resultLoginSuccess.Data.Login.Token)

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	fmt.Println()
	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultSaveCommentToCommentSuccess); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(resultSaveCommentToCommentSuccess.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", resultSaveCommentToCommentSuccess.Errors)
	}

	expectedCommentID := 4
	if resultSaveCommentToCommentSuccess.Data.SaveCommentResponse.ID != expectedCommentID {
		t.Errorf("unexpected comment_id: got %v want %v", resultSaveCommentToCommentSuccess.Data.SaveCommentResponse.ID, expectedCommentID)
	} else {
		fmt.Printf("Comment to comment saved with ID: %v\n", resultSaveCommentToCommentSuccess.Data.SaveCommentResponse.ID)
	}
}

func TestSaveCommentToCommentGraphQLFail1(t *testing.T) {
	query := fmt.Sprintf(`mutation {SaveComment(input: {token:\"%s\", post_id: 1, content: \"hello2\"}) {id created_at} }`, resultLoginFail.Data.Login.Token)

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	fmt.Println()
	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultSaveCommentToCommentFail); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if resultSaveCommentToCommentFail.Data.SaveCommentResponse.ID != 0 {
		t.Errorf("unexpected comment_id: got %v", resultSaveCommentToCommentFail.Data.SaveCommentResponse.ID)
	} else {
		fmt.Printf("Access denied\n")
	}
}

func TestProvideCommentGraphQL1(t *testing.T) {
	query := `mutation { ProvideComment(input: {post_id: 1 parent_id: 1}) { comments {id user_id post_id content created_at parent_id} } }`

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	fmt.Println()
	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultProvideComment); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(resultProvideComment.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", resultProvideComment.Errors)
	}

	if len(resultProvideComment.Data.ProvideCommentResponse.Comment) == 0 {
		fmt.Printf("No comments found\n")
	} else {
		fmt.Printf("Comments received: %v\n",
			resultProvideComment.Data.ProvideCommentResponse.Comment)
	}
}

func TestProvideCommentGraphQ2(t *testing.T) {
	query := `mutation { ProvideComment(input: {post_id: 2 parent_id: 0}) { comments {id user_id post_id content created_at parent_id} } }`

	reqBody := fmt.Sprintf(`{"query":"%s"}`, query)

	resp, err := http.Post("http://localhost:8080/query", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, readErr := ioutil.ReadAll(resp.Body)
		if readErr != nil {
			t.Fatalf("could not read response body: %v", readErr)
		}
		fmt.Printf("Response Body: %s\n", string(respBody))
		t.Fatalf("unexpected status: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	fmt.Println()
	if err := json.NewDecoder(bytes.NewReader(respBody)).Decode(&resultProvideComment); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(resultProvideComment.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", resultProvideComment.Errors)
	}

	if len(resultProvideComment.Data.ProvideCommentResponse.Comment) == 0 {
		fmt.Printf("No comments found\n")
	} else {
		fmt.Printf("Comments received: %v\n",
			resultProvideComment.Data.ProvideCommentResponse.Comment)
	}
}
