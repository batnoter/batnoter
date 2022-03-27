import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { getUserRepos, saveDefaultRepo } from "../api/api";
import { RootState } from "../app/store";

export interface Repo {
  name: string
  visibility: string
  default_branch?: string
}

export enum PreferenceStatus { LOADING, IDLE, FAIL }

interface PreferenceState {
  userRepos: Repo[]
  status: PreferenceStatus
}

const initialState: PreferenceState = {
  userRepos: [],
  status: PreferenceStatus.IDLE
}

export const getUserReposAsync = createAsyncThunk(
  'user/fetchUserRepos',
  async () => {
    const response = await getUserRepos();
    // returned value becomes the `fulfilled` action payload
    return response;
  }
)

export const saveDefaultRepoAsync = createAsyncThunk(
  'user/saveDefaultRepo',
  async (defaultRepo: Repo) => {
    await saveDefaultRepo(defaultRepo);
  }
);

export const preferenceSlice = createSlice({
  name: "preference",
  initialState,
  reducers: {
  },
  extraReducers: (builder) => {
    builder
      .addCase(getUserReposAsync.pending, (state) => {
        state.status = PreferenceStatus.LOADING;
      })
      .addCase(getUserReposAsync.fulfilled, (state, action) => {
        state.status = PreferenceStatus.IDLE;
        state.userRepos = action.payload;
      })
      .addCase(getUserReposAsync.rejected, (state) => {
        state.status = PreferenceStatus.FAIL;
        state.userRepos = [];
      })

      .addCase(saveDefaultRepoAsync.pending, (state) => {
        state.status = PreferenceStatus.LOADING;
      })
      .addCase(saveDefaultRepoAsync.fulfilled, (state) => {
        state.status = PreferenceStatus.IDLE;
      })
      .addCase(saveDefaultRepoAsync.rejected, (state) => {
        state.status = PreferenceStatus.FAIL;
      });
  },
})

export const selectUserRepos = (state: RootState) => state.preference.userRepos;
export const selectPreferenceStatus = (state: RootState) => state.preference.status;
export default preferenceSlice.reducer;
