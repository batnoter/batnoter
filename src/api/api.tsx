import { NotePage, NoteResponsePayload } from "../reducer/noteSlice";
import { Repo } from "../reducer/preferenceSlice";
import { User } from "../reducer/userSlice";

const API_URL = "/api/v1";

const getHeaders = (): HeadersInit => {
  const headers: HeadersInit = new Headers();
  headers.set("Accept", "application/json");
  headers.set("Authorization", "Bearer " + localStorage.getItem("token"));
  headers.set("Content-Type", "application/json");
  return headers;
}

export const getUserProfile = (): Promise<User> => {
  // there is no point in calling profile api without a token, since it will fail anyway due to missing token
  // so in case of page reload etc we call this endpoint only when there is token in local-storage, which implies
  // that user has already logged-in
  return localStorage.getItem("token") ? fetch(`${API_URL}/user/me`, { headers: getHeaders() }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject(await res.json());
    }
    return await res.json();
  }) : Promise.reject("token missing. fetching user profile failed");
}

export const getUserRepos = (): Promise<Repo[]> => {
  return fetch(`${API_URL}/user/preference/repo`, { headers: getHeaders() }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject(await res.json());
    }
    return await res.json();
  })
}

export const autoSetupRepo = (repoName: string): Promise<undefined> => {
  return fetch(`${API_URL}/user/preference/auto/repo?repoName=${repoName}`, {
    method: "POST",
    headers: getHeaders()
  }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject(await res.json());
    }
  })
}

export const saveDefaultRepo = (defaultRepo: Repo): Promise<undefined> => {
  return fetch(`${API_URL}/user/preference/repo`, {
    method: "POST",
    body: JSON.stringify(defaultRepo),
    headers: getHeaders()
  }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject(await res.json());
    }
  })
}

export const searchNotes = (page?: number, path?: string, query?: string): Promise<NotePage> => {
  return fetch(`${API_URL}/search/notes?page=` + (page || 1) + (path ? `path=${path}` : "") + (query ? `query=${query}` : ""),
    { headers: getHeaders() }).then(async (res) => {
      if (!res.ok) {
        return Promise.reject(await res.json());
      }
      return await res.json();
    })
}

export const getNotesTree = (): Promise<NoteResponsePayload[]> => {
  return fetch(`${API_URL}/tree/notes`, {
    headers: getHeaders()
  }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject(await res.json());
    }
    return await res.json();
  })
}

export const getAllNotes = (path: string): Promise<NoteResponsePayload[]> => {
  return fetch(`${API_URL}/notes` + (path && "?path=" + encodeURIComponent(path)), {
    headers: getHeaders()
  }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject(await res.json());
    }
    return await res.json();
  })
}

export const getNote = (path: string): Promise<NoteResponsePayload> => {
  return fetch(`${API_URL}/notes/` + encodeURIComponent(path), {
    headers: getHeaders()
  }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject(await res.json());
    }
    return await res.json();
  })
}

export const saveNote = (path: string, content: string, sha?: string): Promise<NoteResponsePayload> => {
  return fetch(`${API_URL}/notes/` + encodeURIComponent(path), {
    method: "POST",
    body: JSON.stringify({ sha: sha, content: content }),
    headers: getHeaders()
  }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject(await res.json());
    }
    return await res.json();
  })
}

export const deleteNote = (path: string, sha?: string): Promise<undefined> => {
  return fetch(`${API_URL}/notes/` + encodeURIComponent(path), {
    method: "DELETE",
    body: JSON.stringify({ sha: sha }),
    headers: getHeaders()
  }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject(await res.json());
    }
  })
}
