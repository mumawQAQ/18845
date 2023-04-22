from collections import defaultdict

from flask import Flask, request, jsonify, send_from_directory, make_response
import os
import note

app = Flask(__name__)
UPLOADS_DIR = "uploads"
notes = {}
note_id_counter = 1


@app.route("/files", methods=["POST"])
def upload_file():
    if "file" not in request.files:
        return jsonify({"error": "No file provided"}), 400

    file = request.files["file"]
    file.save(os.path.join(UPLOADS_DIR, file.filename))
    return jsonify({"message": "File uploaded successfully", "filename": file.filename}), 200


@app.route("/files", methods=["GET"])
def list_files():
    files = os.listdir(UPLOADS_DIR)
    return jsonify(files), 200


@app.route("/files/<path:filename>", methods=["GET"])
def download_file(filename):
    if not os.path.isfile(os.path.join(UPLOADS_DIR, filename)):
        return jsonify({"error": "File not found"}), 404

    return send_from_directory(UPLOADS_DIR, filename)


@app.route("/files/<path:filename>", methods=["DELETE"])
def delete_file(filename):
    file_path = os.path.join(UPLOADS_DIR, filename)
    if not os.path.isfile(file_path):
        return jsonify({"error": "File not found"}), 404

    os.remove(file_path)
    return jsonify({"message": "File deleted successfully"}), 200

@app.route("/notes", methods=["POST"])
def create_note():
    global note_id_counter
    note_data = request.get_json()

    new_note = note.Note(id=note_id_counter, title=note_data.get("title"), content=note_data.get("content"))
    notes[note_id_counter] = new_note
    note_id_counter += 1

    return jsonify(new_note.to_dict()), 201


@app.route("/notes", methods=["GET"])
def get_all_notes():
    return jsonify([note.to_dict() for note in notes.values()])


@app.route("/notes/<int:note_id>", methods=["GET"])
def get_note_by_id(note_id):
    note = notes.get(note_id)
    if note:
        return jsonify(note.to_dict())
    else:
        return make_response(jsonify({"error": "Note not found"}), 404)


@app.route("/notes/<int:note_id>", methods=["PUT"])
def update_note_by_id(note_id):
    note_data = request.get_json()
    note = notes.get(note_id)

    if not note:
        return make_response(jsonify({"error": "Note not found"}), 404)

    note.title = note_data.get("title", note.title)
    note.content = note_data.get("content", note.content)

    return jsonify(note.to_dict())


@app.route("/notes/<int:note_id>", methods=["DELETE"])
def delete_note_by_id(note_id):
    note = notes.pop(note_id, None)
    if note:
        return jsonify({"message": "Note deleted"})
    else:
        return make_response(jsonify({"error": "Note not found"}), 404)
    
def generate_string(length, base):
    return (base * (length // len(base) + 1))[:length]

def init_notes():
    global notes, note_id_counter
    notes = {}
    for i in range(1, 10001):
        new_note = note.Note(id=i,
                             title=generate_string(10, f"Title {i}"),
                             content=generate_string(100, f"Content {i}"))
        notes[note_id_counter] = new_note
        note_id_counter += 1

init_notes()

if __name__ == "__main__":
    app.run()
