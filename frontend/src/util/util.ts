import { SerializedError } from "@reduxjs/toolkit";
import { ShowFn } from "mui-modal-provider/dist/types";
import ConfirmDialog from "../components/ConfirmDialog";


const REPLACE_EXT_REGEX = /(\.md)$/i;
const EXT = '.md';
const BACKEND_ERROR_CODES = ['internal_server_error', 'validation_failed'];
const UNKNOWN_ERR_MSG = "Something went wrong. Please try again!"

export function getTitleFromFilename(filename: string): string {
  return filename.replace(REPLACE_EXT_REGEX, '');
}

export function getFilenameFromTitle(title: string): string {
  return title + EXT;
}

export function getPathWithoutExt(path: string): string {
  return path.replace(REPLACE_EXT_REGEX, '');
}

export function getDecodedPath(path: string | null): string {
  if (path == null) {
    return "";
  }
  const decodedPath = decodeURIComponent(path || "");
  if (decodedPath === '/') {
    return "";
  }
  return decodedPath;
}

export function appendPath(parentPath: string, path: string) {
  if (parentPath === "") {
    return path;
  }
  if (path === "") {
    return parentPath;
  }
  return parentPath + '/' + path;
}

export function isFilePath(path: string): boolean {
  return path.endsWith(EXT);
}

export function splitPath(path: string): string[] {
  // split the path and return the array ignoring any blank elements
  return path.split('/').filter(p => p);
}

export function confirmDeleteNote(showModal: ShowFn, onConfirm: () => void) {
  showModal(ConfirmDialog, {
    desc: 'Are you sure you want to delete this note?',
    onConfirm: onConfirm
  });
}

export function getSanitizedErrorMessage(error: SerializedError): string {
  // Validate error. Since we don't want to show programming errors to users.
  if (error.code != null && error.message != null && BACKEND_ERROR_CODES.includes(error.code)) {
    return error.message;
  }
  return UNKNOWN_ERR_MSG;
}
