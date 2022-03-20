import { Action, configureStore, ThunkAction } from '@reduxjs/toolkit';
import noteReducer from '../reducer/note/noteSlice';
import userReducer from '../reducer/user/userSlice';

export const store = configureStore({
  reducer: {
    user: userReducer,
    notes: noteReducer
  },
});

export type AppDispatch = typeof store.dispatch;
export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<
  ReturnType,
  RootState,
  unknown,
  Action<string>
>;
