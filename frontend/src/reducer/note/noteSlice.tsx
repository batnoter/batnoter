import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { deleteNote, getAllNotes, getNoteByID, saveNote } from "../../api/api";
import { RootState } from "../../app/store";

export interface Note {
  id?: number
  created_at?: string
  updated_at?: string
  title: string
  content: string
}

export enum NoteStatus { LOADING, IDLE, FAIL }

interface NoteState {
  list: Note[]
  current: Note | null
  status: NoteStatus
}

const initialState: NoteState = {
  list: [],
  current: null,
  status: NoteStatus.IDLE
}

export const getAllNotesAsync = createAsyncThunk(
  'note/fetchNotes',
  async () => {
    const response = await getAllNotes();
    return response;
  }
);

export const getNoteByIDAsync = createAsyncThunk(
  'note/fetchNoteByID',
  async (id: string) => {
    const response = await getNoteByID(id);
    return response;
  }
);

export const saveNoteAsync = createAsyncThunk(
  'note/saveNote',
  async (note: Note) => {
    const response = await saveNote(note);
    return response;
  }
);

export const deleteNoteAsync = createAsyncThunk(
  'note/deleteNote',
  async (id: number) => {
    await deleteNote(id);
    return id;
  }
);

export const noteSlice = createSlice({
  name: "notes",
  initialState,
  reducers: {
  },
  extraReducers: (builder) => {
    builder
      .addCase(getAllNotesAsync.pending, (state) => {
        state.status = NoteStatus.LOADING;
      })
      .addCase(getAllNotesAsync.fulfilled, (state, action) => {
        state.list = action.payload;
        state.status = NoteStatus.IDLE;
      })
      .addCase(getAllNotesAsync.rejected, (state) => {
        state.list = [];
        state.status = NoteStatus.FAIL;
      })

      .addCase(getNoteByIDAsync.pending, (state) => {
        state.current = null
        state.status = NoteStatus.LOADING;
      })
      .addCase(getNoteByIDAsync.fulfilled, (state, action) => {
        state.current = action.payload
        state.status = NoteStatus.IDLE;
      })
      .addCase(getNoteByIDAsync.rejected, (state) => {
        state.status = NoteStatus.FAIL;
      })

      .addCase(saveNoteAsync.pending, (state) => {
        state.status = NoteStatus.LOADING;
      })
      .addCase(saveNoteAsync.fulfilled, (state, action) => {
        state.list = state.list.filter(n => n.id !== action.payload)
        state.list.push(action.payload)
        state.status = NoteStatus.IDLE;
      })
      .addCase(saveNoteAsync.rejected, (state) => {
        state.status = NoteStatus.FAIL;
      })

      .addCase(deleteNoteAsync.pending, (state) => {
        state.status = NoteStatus.LOADING;
      })
      .addCase(deleteNoteAsync.fulfilled, (state, action) => {
        state.list = state.list.filter(n => n.id !== action.payload)
        state.status = NoteStatus.IDLE;
      })
      .addCase(deleteNoteAsync.rejected, (state) => {
        state.status = NoteStatus.FAIL;
      });
  },
})
export const selectCurrentNote = (state: RootState) => state.notes.current;
export const selectNotes = (state: RootState) => state.notes.list;
export const selectNoteStatus = (state: RootState) => state.notes.status;
export default noteSlice.reducer;