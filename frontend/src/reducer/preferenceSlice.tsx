import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { autoSetupRepo, getUserRepos, saveDefaultRepo } from "../api/api";
import { RootState } from "../app/store";
import { APIStatus, APIStatusType } from "./common";

export interface Repo {
  name: string
  visibility: string
  default_branch?: string
}

interface PreferenceState {
  userRepos: Repo[]
  status: APIStatus
}

const initialState: PreferenceState = {
  userRepos: [],
  status: {
    getUserReposAsync: APIStatusType.IDLE,
    autoSetupRepoAsync: APIStatusType.IDLE,
    saveDefaultRepoAsync: APIStatusType.IDLE,
  }
}

export const getUserReposAsync = createAsyncThunk(
  'user/fetchUserRepos',
  async () => {
    const response = await getUserRepos();
    // returned value becomes the `fulfilled` action payload
    return response;
  }
)

export const autoSetupRepoAsync = createAsyncThunk(
  'user/autoSetupRepo',
  async (repoName: string) => {
    await autoSetupRepo(repoName);
  }
);

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
        state.status.getUserReposAsync = APIStatusType.LOADING;
      })
      .addCase(getUserReposAsync.fulfilled, (state, action) => {
        state.status.getUserReposAsync = APIStatusType.IDLE;
        state.userRepos = action.payload as Repo[];
      })
      .addCase(getUserReposAsync.rejected, (state) => {
        state.status.getUserReposAsync = APIStatusType.FAIL;
        state.userRepos = [];
      })

      .addCase(autoSetupRepoAsync.pending, (state) => {
        state.status.autoSetupRepoAsync = APIStatusType.LOADING;
      })
      .addCase(autoSetupRepoAsync.fulfilled, (state) => {
        state.status.autoSetupRepoAsync = APIStatusType.IDLE;
      })
      .addCase(autoSetupRepoAsync.rejected, (state) => {
        state.status.autoSetupRepoAsync = APIStatusType.FAIL;
      })

      .addCase(saveDefaultRepoAsync.pending, (state) => {
        state.status.saveDefaultRepoAsync = APIStatusType.LOADING;
      })
      .addCase(saveDefaultRepoAsync.fulfilled, (state) => {
        state.status.saveDefaultRepoAsync = APIStatusType.IDLE;
      })
      .addCase(saveDefaultRepoAsync.rejected, (state) => {
        state.status.saveDefaultRepoAsync = APIStatusType.FAIL;
      });
  },
})

export const selectUserRepos = (state: RootState): Repo[] => state.preference.userRepos;
export const selectPreferenceAPIStatus = (state: RootState): APIStatus => state.preference.status;
export default preferenceSlice.reducer;
