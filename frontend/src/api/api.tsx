import { Note } from "../reducer/note/noteSlice";

// when we start the application using "npm start" then process.env.NODE_ENV will be automatically set to "development"
const API_URL = "/api/v1";

const getHeaders = () => {
  return {
    Accept: "application/json",
    Authorization: "Bearer " + localStorage.getItem("token"),
    "Content-Type": "application/json"
  };
};

const catchError = (error: string) => {
  console.log(error);
  return Promise.reject(error)
};

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
};

export const getAllNotes = () => {
  return fetch(`${API_URL}/note`, { headers: getHeaders() }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject("fetching notes failed")
    }
    return await res.json();
  }).catch(catchError)
}

export const getNoteByID = (id: string) => {
  return fetch(`${API_URL}/note/${id}`, { headers: getHeaders() }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject("fetching note failed")
    }
    return await res.json();
  }).catch(catchError)
}

export const saveNote = (note: Note) => {
  return fetch(`${API_URL}/note/`, {
    method: "POST",
    body: JSON.stringify(note),
    headers: getHeaders()
  }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject("saving note failed")
    }
    return await res.json();
  }).catch(catchError)
}

export const deleteNote = (id: number) => {
  return fetch(`${API_URL}/note/${id}`, { method: "DELETE", headers: getHeaders() }).then(async (res) => {
    if (!res.ok) {
      return Promise.reject("deleting note failed")
    }
  }).catch(catchError)
}