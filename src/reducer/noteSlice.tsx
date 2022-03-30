import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { deleteNote, getAllNotes, getNote, getNotesTree, saveNote, searchNotes } from "../api/api";
import { RootState } from "../app/store";

export interface SearchParams {
  page?: number
  path?: string
  query?: string
}

export interface Tree {
  [key: string]: any // type for unknown keys.
  children?: Tree[] // type for a known property.
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
  tree: Tree
  current: Note | null
  status: NoteStatus
}

const initialState: NoteState = {
  page: {
    total: 1,
    notes: []
  },
  tree: {
    name: "root",
    path: "",
    cached: false
  },
  current: null,
  status: NoteStatus.IDLE
}

export const searchNotesAsync = createAsyncThunk(
  'note/searchNotes',
  async (params?: SearchParams) => {
    const response = await searchNotes(params?.page, params?.path, params?.query);
    return response;
  }
);

export const getNotesTreeAsync = createAsyncThunk(
  'note/fetchNotesTree',
  async () => {
    const response = await getNotesTree();
    return response;
  }
);

export const getNotesAsync = createAsyncThunk(
  'note/fetchNotes',
  async (path: string) => {
    const response = await getAllNotes(path);
    return response;
  }, {
  condition: (path, { getState }) => {
    const state = getState() as RootState
    const node = TreeUtil.searchNode(state.notes.tree, path)
    const hasFiles = !!(node?.children && node.children.find(o => !o.is_dir))
    return !node?.cached && hasFiles
  }
}
);

export const getNoteAsync = createAsyncThunk(
  'note/fetchNote',
  async (path: string) => {
    const response = await getNote(path);
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
        const tree = TreeUtil.parse(state.tree, state.page.notes, true);
        state.tree = tree;
        state.status = NoteStatus.IDLE;
      })
      .addCase(searchNotesAsync.rejected, (state) => {
        state.page = initialState.page;
        state.status = NoteStatus.FAIL;
      })

      .addCase(getNotesTreeAsync.pending, (state) => {
        state.status = NoteStatus.LOADING;
      })
      .addCase(getNotesTreeAsync.fulfilled, (state, action) => {
        state.page.notes = action.payload as Note[];
        const tree = TreeUtil.parse(initialState.tree, state.page.notes, false);
        state.tree = tree;
        state.status = NoteStatus.IDLE;
      })
      .addCase(getNotesTreeAsync.rejected, (state) => {
        state.page.notes = initialState.page.notes;
        state.status = NoteStatus.FAIL;
      })

      .addCase(getNotesAsync.pending, (state) => {
        state.status = NoteStatus.LOADING;
      })
      .addCase(getNotesAsync.fulfilled, (state, action) => {
        state.page.notes = action.payload as Note[];
        const tree = TreeUtil.parse(state.tree, state.page.notes, true);
        state.tree = tree;
        state.status = NoteStatus.IDLE;
      })
      .addCase(getNotesAsync.rejected, (state) => {
        state.page.notes = initialState.page.notes;
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
        state.page.notes = state.page.notes.filter(n => n.path !== action.payload.path)
        TreeUtil.deleteNode(state.tree, action.payload.path)
        state.status = NoteStatus.IDLE;
      })
      .addCase(deleteNoteAsync.rejected, (state) => {
        state.status = NoteStatus.FAIL;
      });
  },
})

export class TreeUtil {
  static parse(seedTree: Tree, notes: Note[], cache: boolean): Tree {
    const tree: Tree = notes.reduce((r, n) => {
      const path = n.path.split('/')
      const fileName = path.pop() || ""
      const final = path.reduce((o, name) => {
        let temp = (o.children = o.children || []).find(q => q.name === name);
        if (!temp) o.children.push(temp = {
          name,
          path: o.path ? o.path + '/' + name : name,
          is_dir: true
        });
        temp.cached = cache;
        return temp;
      }, r);

      const file = { ...n, name: fileName }
      final.children = final.children || [];
      const index = final.children.findIndex(o => o.path === n.path);
      index > -1 && (final.children[index] = file) || final.children.push(file);
      final.cached = cache
      return r;
    }, { ...seedTree });

    return tree
  }

  static searchNode(root: Tree, path: string): Tree | null {
    if (root.path == path) {
      return root;
    }

    if (root.children != null) {
      let result = null;
      for (let i = 0; result == null && i < root.children.length; i++) {
        result = TreeUtil.searchNode(root.children[i], path);
      }
      return result;
    }
    return null;
  }

  static deleteNode(root: Tree, path: string) {
    if (!root.children) {
      return
    }
    for (let i = 0; i < root.children.length; i++) {
      if (root.children[i].path == path) {
        root.children.splice(i, 1)
        break
      }
      TreeUtil.deleteNode(root.children[i], path)
    }
  }
}


export const selectCurrentNote = (state: RootState): Note | null => state.notes.current;
export const selectNotesPage = (state: RootState): NotePage => state.notes.page;
export const selectNotesTree = (state: RootState): Tree => state.notes.tree;
export const selectNoteStatus = (state: RootState): NoteStatus => state.notes.status;
export default noteSlice.reducer;
