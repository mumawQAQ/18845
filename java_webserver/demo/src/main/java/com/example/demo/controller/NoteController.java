package com.example.demo.controller;

import com.example.demo.model.Note;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Collection;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;

@RestController
@RequestMapping("/notes")
public class NoteController {

    private final Map<Integer, Note> notes = new HashMap<>();
    private final AtomicInteger noteIdCounter = new AtomicInteger(1);

    private String generateString(int length, String base) {
        return String.format("%1$-" + length + "s", base).substring(0, length);
    }

    @PostMapping("/init_notes")
    @ResponseStatus(HttpStatus.OK)
    public String initNotes() {
        notes.clear();
        for (int i = 1; i <= 10000; i++) {
            int id = noteIdCounter.getAndIncrement();
            notes.put(id, new Note(id,
                                   generateString(10, "Title " + i),
                                   generateString(100, "Content " + i)));
        }
        return "10,000 notes initialized";
    }

    @PostMapping
    public ResponseEntity<Note> createNote(@RequestBody Note note) {
        int noteId = noteIdCounter.getAndIncrement();
        Note newNote = new Note(noteId, note.getTitle(), note.getContent());
        notes.put(noteId, newNote);
        return new ResponseEntity<>(newNote, HttpStatus.CREATED);
    }

    @GetMapping
    public Collection<Note> getAllNotes() {
        return notes.values();
    }

    @GetMapping("/{id}")
    public ResponseEntity<Note> getNoteById(@PathVariable int id) {
        Note note = notes.get(id);
        if (note != null) {
            return ResponseEntity.ok(note);
        } else {
            return ResponseEntity.notFound().build();
        }
    }

    @PutMapping("/{id}")
    public ResponseEntity<Note> updateNoteById(@PathVariable int id, @RequestBody Note updatedNote) {
        Note note = notes.get(id);
        if (note != null) {
            note.setTitle(updatedNote.getTitle());
            note.setContent(updatedNote.getContent());
            return ResponseEntity.ok(note);
        } else {
            return ResponseEntity.notFound().build();
        }
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteNoteById(@PathVariable int id) {
        if (notes.containsKey(id)) {
            notes.remove(id);
            return ResponseEntity.ok().build();
        } else {
            return ResponseEntity.notFound().build();
        }
    }
}
