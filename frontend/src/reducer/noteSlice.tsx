import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { deleteNote, getAllNotes, getNote, getNotesTree, saveNote, searchNotes } from "../api/api";
import { RootState } from "../app/store";

export interface SearchParams {
  page?: number
  path?: string
  query?: string
}

export interface TreeNode {
  name: string
  sha?: string
  path: string
  content?: string
  size?: number
  is_dir: boolean
  cached: boolean
  children?: TreeNode[]
}

export interface NoteResponsePayload {
  sha: string
  path: string
  content: string
  size: number
  is_dir: boolean
}

export interface NotePage {
  total: number
  notes: NoteResponsePayload[]
}

export enum NoteStatus { LOADING, IDLE, FAIL }

interface NoteState {
  page: NotePage
  tree: TreeNode
  current: NoteResponsePayload | null
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
    cached: false,
    is_dir: true
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
    const response = await getNotesTree() as NoteResponsePayload[];
    return response;
  }
);

export const getNotesAsync = createAsyncThunk(
  'note/fetchNotes',
  async (path: string) => {
    const response = await getAllNotes(path) as NoteResponsePayload[];
    return response;
  }, {
  condition: (path, { getState }) => {
    const state = getState() as RootState;
    const node = TreeUtil.searchNode(state.notes.tree, path);
    const hasFiles = !!(node?.children && node.children.find(o => !o.is_dir));
    return !node?.cached && hasFiles;
  }
}
);

export const getNoteAsync = createAsyncThunk(
  'note/fetchNote',
  async (path: string) => {
    const response = await getNote(path) as NoteResponsePayload;
    return response;
  }
);

export const saveNoteAsync = createAsyncThunk(
  'note/saveNote',
  async ({ path, content, sha }: { path: string, content: string, sha?: string }) => {
    const response = await saveNote(path, content, sha) as NoteResponsePayload;
    return {
      ...response,
      content: content
    };
  }
);

export const deleteNoteAsync = createAsyncThunk(
  'note/deleteNote',
  async (note: TreeNode) => {
    await deleteNote(note.path, note.sha);
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
        state.page.notes = action.payload;
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
        state.page.notes = action.payload;
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
        state.current = action.payload;
        const tree = TreeUtil.parse(state.tree, [action.payload]);
        state.tree = tree;
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
        const tree = TreeUtil.parse(state.tree, [action.payload]);
        state.tree = tree;
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
  static parse(seedTree: TreeNode, notes: NoteResponsePayload[], cache?: boolean): TreeNode {
    const tree: TreeNode = notes.reduce((r, n) => {
      const pathArray = n.path.split('/');
      const fileName = pathArray.pop() || "";
      const final = pathArray.reduce((o, name) => {
        let temp = (o.children = o.children || []).find(q => q.name === name);
        if (!temp) o.children.push(temp = {
          name,
          path: o.path ? o.path + '/' + name : name,
          is_dir: true,
          cached: !!cache
        });
        cache != null && (temp.cached = cache);
        o.children.sort((a, b) => (Number(b.is_dir) - Number(a.is_dir)) || a.path.localeCompare(b.path))
        return temp;
      }, r);

      const file = { ...n, name: fileName, cached: !!n.content }
      final.children = final.children || [];
      const index = final.children.findIndex(o => o.path === n.path);
      index > -1 && (final.children[index] = file) || final.children.push(file);
      final.children.sort((a, b) => (Number(b.is_dir) - Number(a.is_dir)) || a.path.localeCompare(b.path))
      cache != null && (final.cached = cache)
      return r;
    }, { ...seedTree });

    return tree;
  }

  static searchNode(root: TreeNode, path: string): TreeNode | null {
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

  static deleteNode(root: TreeNode, path: string) {
    if (!root.children) {
      return;
    }
    for (let i = 0; i < root.children.length; i++) {
      if (root.children[i].path == path) {
        root.children.splice(i, 1);
        break;
      }
      TreeUtil.deleteNode(root.children[i], path);
    }
  }

  static getChildDirs(tree: TreeNode, path: string): string[] {
    const node = TreeUtil.searchNode(tree, path)
    if (!node?.children) {
      return [];
    }
    return node.children.filter(c => c.is_dir).map(c => c.name);
  }
}


export const selectCurrentNote = (state: RootState): NoteResponsePayload | null => state.notes.current;
export const selectNotesPage = (state: RootState): NotePage => state.notes.page;
export const selectNotesTree = (state: RootState): TreeNode => state.notes.tree;
export const selectNoteStatus = (state: RootState): NoteStatus => state.notes.status;
export default noteSlice.reducer;
