import { Action, configureStore, ThunkAction } from '@reduxjs/toolkit';
import noteReducer from '../reducer/noteSlice';
import preferenceReducer from '../reducer/preferenceSlice';
import userReducer from '../reducer/userSlice';

export const store = configureStore({
  reducer: {
    user: userReducer,
    notes: noteReducer,
    preference:  preferenceReducer
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
