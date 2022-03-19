import { createAsyncThunk, createSlice, PayloadAction } from "@reduxjs/toolkit";
import { getUserProfile } from "../../api/api";
import { RootState } from "../../app/store";

export interface User {
    email: string
    name: string
    location: string
    avatar_url: string
}

interface UserState {
    value: User | null
    status: 'loading' | 'idle' | 'fail'
}

const initialState: UserState = {
    value: null,
    status: 'idle'
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
        login: (state, action: PayloadAction<User | null>) => {
            state.value = action.payload
        },
        logout: (state) => {
            state.value = null
            localStorage.removeItem("token")
        }
    },
    extraReducers: (builder) => {
        builder
            .addCase(getUserProfileAsync.pending, (state) => {
                state.status = 'loading';
            })
            .addCase(getUserProfileAsync.fulfilled, (state, action) => {
                state.status = 'idle';
                state.value = action.payload;
            })
            .addCase(getUserProfileAsync.rejected, (state) => {
                state.status = 'fail';
                state.value = null;
                localStorage.removeItem("token")
            });
    },
})

export const { login, logout } = userSlice.actions
export const selectUser = (state: RootState) => state.user.value
export default userSlice.reducer