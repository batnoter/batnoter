import { Masonry } from '@mui/lab';
import { Container } from '@mui/material';
import React, { useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '../app/hooks';
import { deleteNoteAsync, getAllNotesAsync, selectNotes } from '../reducer/note/noteSlice';
import NoteCard from './NoteCard';


const Finder = () => {
  const dispatch = useAppDispatch();

  useEffect(() => {
    dispatch(getAllNotesAsync())
  }, []);

  const notes = useAppSelector(selectNotes);

  const handleDelete = (id: number) => {
    dispatch(deleteNoteAsync(id))
  }

  return (
    <Container>
      <Masonry columns={{ xs: 1, sm: 3, xl: 4 }} spacing={2}>
        {notes.map(note => (
          <div key={note.id}>
            <NoteCard note={note} handleDelete={handleDelete} />
          </div>
        ))}
      </Masonry>
    </Container>
  )
}

export default Finder