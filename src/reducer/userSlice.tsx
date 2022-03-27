import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { getUserProfile } from "../api/api";
import { RootState } from "../app/store";

export interface User {
  email: string
  name: string
  location: string
  avatar_url: string
  default_repo? : {
    name: string,
    visibility: string,
    default_branch: string
  }
}

export enum UserStatus { LOADING, IDLE, FAIL }

interface UserState {
  value: User | null
  status: UserStatus
}

const initialState: UserState = {
  value: null,
  status: UserStatus.IDLE
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
      state.status = UserStatus.LOADING;
    },
    userLogout: (state) => {
      state.value = null;
      localStorage.removeItem("token");
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(getUserProfileAsync.pending, (state) => {
        state.status = UserStatus.LOADING;
      })
      .addCase(getUserProfileAsync.fulfilled, (state, action) => {
        state.status = UserStatus.IDLE;
        state.value = action.payload;
      })
      .addCase(getUserProfileAsync.rejected, (state) => {
        state.status = UserStatus.FAIL;
        state.value = null;
        localStorage.removeItem("token")
      });
  },
})

export const { userLoading, userLogout } = userSlice.actions;
export const selectUser = (state: RootState) => state.user.value;
export const selectUserStatus = (state: RootState) => state.user.status;
export default userSlice.reducer;
