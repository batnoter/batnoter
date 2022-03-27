import { Note } from "../reducer/noteSlice";
import { Repo } from "../reducer/preferenceSlice";

// when we start the application using "npm start" then process.env.NODE_ENV will be automatically set to "development"
const API_URL = "/api/v1";

const getHeaders = () => {
  return {
    Accept: "application/json",
    Authorization: "Bearer " + localStorage.getItem("token"),
    "Content-Type": "application/json"
  }
}

const catchError = (error: string) => {
  console.log(error);
  return Promise.reject(error)
}

export const getUserProfile = () => {
  // There is no point in calling profile api without a token, since it will fail anyway due to missing token
  // So in case of page reload etc we call this endpoint only when there is token in local-storage, which implies
  // that user has already logged-in
  return localStorage.getItem("token") ? fetch(`${API_URL}/user/me`, { headers: getHeaders() }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject("fetching user profile failed")
    }
    return await res.json();
  }).catch(catchError) : Promise.reject("token missing. fetching user profile failed");
}

export const getUserRepos = () => {
  return fetch(`${API_URL}/user/preference/repo`, { headers: getHeaders() }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject("fetching user's github repos failed")
    }
    return await res.json();
  }).catch(catchError);
}

export const saveDefaultRepo = (defaultRepo: Repo) => {
  return fetch(`${API_URL}/user/preference/repo`, {
    method: "POST",
    body: JSON.stringify(defaultRepo),
    headers: getHeaders()
  }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject("saving user's default repo failed")
    }
  }).catch(catchError)
}

export const searchNotes = (page?: number, path?: string, query?: string) => {
  return fetch(`${API_URL}/note?page=` + (page || 1) + (path ? `path=${path}` : "") + (query ? `query=${query}` : ""),
    { headers: getHeaders() }).then(async (res) => {
      if (!res.ok) {
        return Promise.reject("fetching notes failed")
      }
      return await res.json();
    }).catch(catchError)
}

export const getNote = (path: string) => {
  return fetch(`${API_URL}/note/` + encodeURIComponent(path), {
    headers: getHeaders()
  }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject("fetching note failed")
    }
    return await res.json();
  }).catch(catchError)
}

export const saveNote = (path: string, content: string, sha?: string) => {
  return fetch(`${API_URL}/note/` + encodeURIComponent(path), {
    method: "POST",
    body: JSON.stringify({ sha: sha, content: content }),
    headers: getHeaders()
  }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject("saving note failed")
    }
    return await res.json();
  }).catch(catchError)
}

export const deleteNote = (note: Note) => {
  return fetch(`${API_URL}/note/` + encodeURIComponent(note.path), {
    method: "DELETE",
    body: JSON.stringify({ sha: note.sha }),
    headers: getHeaders()
  }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject("deleting note failed")
    }
  }).catch(catchError)
}