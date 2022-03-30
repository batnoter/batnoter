import { Masonry } from '@mui/lab';
import { Container } from '@mui/material';
import React, { ReactElement, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../app/hooks';
import { deleteNoteAsync, getNotesAsync, Note, selectNotesTree, Tree, TreeUtil } from '../reducer/noteSlice';
import NoteCard from './NoteCard';


const Finder = (): ReactElement => {
  const dispatch = useAppDispatch();
  const tree = useAppSelector(selectNotesTree)
  const [searchParams] = useSearchParams();
  let path = decodeURIComponent(searchParams.get('path') || "")
  path = path === "/" ? "" : path
  useEffect(() => {
    dispatch(getNotesAsync(path))
  }, [path, tree])

  const handleDelete = (note: Note) => {
    dispatch(deleteNoteAsync(note))
  }

  const getChildren = (path: string): Tree[] | undefined => {
    const node = TreeUtil.searchNode(tree, path)
    if (node?.cached) {
      return node.children
    }
  }

  const notes = getChildren(path) || [] as Tree[]

  return (
    <Container>
      <Masonry columns={{ xs: 1, md: 3, xl: 4 }} spacing={2}>
        {notes.filter(n => !n.is_dir).map(note => (
          <div key={note.path}>
            <NoteCard note={note as Note} handleDelete={handleDelete} />
          </div>
        ))}
      </Masonry>
    </Container>
  )
}

export default Finder