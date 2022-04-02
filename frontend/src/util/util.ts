

const REPLACE_EXT_REGEX = /(\.md)$/i;
const EXT = '.ms';

export function getTitleFromFilename(filename: string): string {
  return filename.replace(REPLACE_EXT_REGEX, '');
}

export function getFilenameFromTitle(title: string): string {
  return title + EXT;
}