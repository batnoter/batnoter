import { Note, NotePage } from "../reducer/noteSlice";
import { Repo } from "../reducer/preferenceSlice";
import { User } from "../reducer/userSlice";

// when we start the application using "npm start" then process.env.NODE_ENV will be automatically set to "development"
const API_URL = "/api/v1";

const getHeaders = (): HeadersInit => {
  const headers: HeadersInit = new Headers();
  headers.set("Accept", "application/json")
  headers.set("Authorization", "Bearer " + localStorage.getItem("token"))
  headers.set("Content-Type", "application/json")
  return headers
}

const catchError = (error: string): Promise<string> => {
  console.log(error);
  return Promise.reject(error)
}

export const getUserProfile = (): Promise<User | string> => {
  // there is no point in calling profile api without a token, since it will fail anyway due to missing token
  // so in case of page reload etc we call this endpoint only when there is token in local-storage, which implies
  // that user has already logged-in
  return localStorage.getItem("token") ? fetch(`${API_URL}/user/me`, { headers: getHeaders() }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject("fetching user profile failed")
    }
    return await res.json();
  }).catch(catchError) : Promise.reject("token missing. fetching user profile failed");
}

export const getUserRepos = (): Promise<Repo[] | string> => {
  return fetch(`${API_URL}/user/preference/repo`, { headers: getHeaders() }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject("fetching user's github repos failed")
    }
    return await res.json();
  }).catch(catchError);
}

export const saveDefaultRepo = (defaultRepo: Repo): Promise<undefined | string> => {
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

export const searchNotes = (page?: number, path?: string, query?: string): Promise<NotePage | string> => {
  return fetch(`${API_URL}/search/notes?page=` + (page || 1) + (path ? `path=${path}` : "") + (query ? `query=${query}` : ""),
    { headers: getHeaders() }).then(async (res) => {
      if (!res.ok) {
        return Promise.reject("fetching notes failed")
      }
      return await res.json();
    }).catch(catchError)
}

export const getNotesTree = (): Promise<Note[] | string> => {
  return fetch(`${API_URL}/tree/notes`, {
    headers: getHeaders()
  }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject("fetching notes tree failed")
    }
    return await res.json();
  }).catch(catchError)
}

export const getAllNotes = (path: string): Promise<Note[] | string> => {
  return fetch(`${API_URL}/notes` + (path && "?path=" + encodeURIComponent(path)), {
    headers: getHeaders()
  }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject("fetching notes failed")
    }
    return await res.json();
  }).catch(catchError)
}

export const getNote = (path: string): Promise<Note | string> => {
  return fetch(`${API_URL}/notes/` + encodeURIComponent(path), {
    headers: getHeaders()
  }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject("fetching note failed")
    }
    return await res.json();
  }).catch(catchError)
}

export const saveNote = (path: string, content: string, sha?: string): Promise<Note | string> => {
  return fetch(`${API_URL}/notes/` + encodeURIComponent(path), {
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

export const deleteNote = (note: Note): Promise<undefined | string> => {
  return fetch(`${API_URL}/notes/` + encodeURIComponent(note.path), {
    method: "DELETE",
    body: JSON.stringify({ sha: note.sha }),
    headers: getHeaders()
  }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject("deleting note failed")
    }
  }).catch(catchError)
}