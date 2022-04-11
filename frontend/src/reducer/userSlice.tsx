import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { getUserProfile } from "../api/api";
import { RootState } from "../app/store";
import { APIStatusType } from "./common";

export interface User {
  email: string
  name: string
  location: string
  avatar_url: string
  default_repo?: {
    name: string,
    visibility: string,
    default_branch: string
  }
}

interface UserState {
  value: User | null
  status: APIStatusType
}

const initialState: UserState = {
  value: null,
  status: APIStatusType.IDLE
}

export const getUserProfileAsync = createAsyncThunk(
  'user/fetchUser',
  async () => {
    const response = await getUserProfile();
    // returned value becomes the `fulfilled` action payload
    return response;
  }
);

export const userSlice = createSlice({
  name: "user",
  initialState,
  reducers: {
    userLoading: (state) => {
      state.status = APIStatusType.LOADING;
    },
    userLogout: (state) => {
      state.value = null;
      localStorage.removeItem("token");
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(getUserProfileAsync.pending, (state) => {
        state.status = APIStatusType.LOADING;
      })
      .addCase(getUserProfileAsync.fulfilled, (state, action) => {
        state.status = APIStatusType.IDLE;
        state.value = action.payload as User;
      })
      .addCase(getUserProfileAsync.rejected, (state) => {
        state.status = APIStatusType.FAIL;
        state.value = null;
        localStorage.removeItem("token")
      });
  },
})

export const { userLoading, userLogout } = userSlice.actions;
export const selectUser = (state: RootState): User | null => state.user.value;
export const selectUserAPIStatus = (state: RootState): APIStatusType => state.user.status;
export default userSlice.reducer;
