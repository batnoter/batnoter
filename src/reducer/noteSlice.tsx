import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { deleteNote, getNote, saveNote, searchNotes } from "../api/api";
import { RootState } from "../app/store";

export interface SearchParams {
  page?: number
  path?: string
  query?: string
}

export interface Note {
  sha: string
  path: string
  content: string
  size: number
  is_dir: boolean
}

export interface NotePage {
  total: number
  notes: Note[]
}

export enum NoteStatus { LOADING, IDLE, FAIL }

interface NoteState {
  page: NotePage
  current: Note | null
  status: NoteStatus
}

const initialState: NoteState = {
  page: {
    total: 0,
    notes: []
  },
  current: null,
  status: NoteStatus.IDLE
}

export const searchNotesAsync = createAsyncThunk(
  'note/fetchNotes',
  async (params?: SearchParams) => {
    const response = await searchNotes(params?.page, params?.path, params?.query);
    return response;
  }
);

export const getNoteAsync = createAsyncThunk(
  'note/fetchNote',
  async (id: string) => {
    const response = await getNote(id);
    return response;
  }
);

export const saveNoteAsync = createAsyncThunk(
  'note/saveNote',
  async ({ path, content, sha }: { path: string, content: string, sha?: string }) => {
    const response = await saveNote(path, content, sha) as Note;
    return {
      ...response,
      content: content
    };
  }
);

export const deleteNoteAsync = createAsyncThunk(
  'note/deleteNote',
  async (note: Note) => {
    await deleteNote(note);
    return note;
  }
);

export const noteSlice = createSlice({
  name: "notes",
  initialState,
  reducers: {
  },
  extraReducers: (builder) => {
    builder
      .addCase(searchNotesAsync.pending, (state) => {
        state.status = NoteStatus.LOADING;
      })
      .addCase(searchNotesAsync.fulfilled, (state, action) => {
        state.page = action.payload as NotePage;
        state.status = NoteStatus.IDLE;
      })
      .addCase(searchNotesAsync.rejected, (state) => {
        state.page = initialState.page;
        state.status = NoteStatus.FAIL;
      })

      .addCase(getNoteAsync.pending, (state) => {
        state.current = null
        state.status = NoteStatus.LOADING;
      })
      .addCase(getNoteAsync.fulfilled, (state, action) => {
        state.current = action.payload as Note;
        state.status = NoteStatus.IDLE;
      })
      .addCase(getNoteAsync.rejected, (state) => {
        state.status = NoteStatus.FAIL;
      })

      .addCase(saveNoteAsync.pending, (state) => {
        state.status = NoteStatus.LOADING;
      })
      .addCase(saveNoteAsync.fulfilled, (state, action) => {
        state.page.notes = state.page.notes.filter(n => n.sha !== action.payload.sha)
        state.page.notes.push(action.payload)
        state.status = NoteStatus.IDLE;
      })
      .addCase(saveNoteAsync.rejected, (state) => {
        state.status = NoteStatus.FAIL;
      })

      .addCase(deleteNoteAsync.pending, (state) => {
        state.status = NoteStatus.LOADING;
      })
      .addCase(deleteNoteAsync.fulfilled, (state, action) => {
        state.page.notes = state.page.notes.filter(n => n.sha !== action.payload.sha)
        state.status = NoteStatus.IDLE;
      })
      .addCase(deleteNoteAsync.rejected, (state) => {
        state.status = NoteStatus.FAIL;
      });
  },
})
export const selectCurrentNote = (state: RootState): Note | null => state.notes.current;
export const selectNotesPage = (state: RootState): NotePage => state.notes.page;
export const selectNoteStatus = (state: RootState): NoteStatus => state.notes.status;
export default noteSlice.reducer;