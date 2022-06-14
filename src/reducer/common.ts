export enum APIStatusType { LOADING, IDLE, FAIL }

export interface APIStatus {
  [asyncName: string]: APIStatusType
}

export type ThemeMode = 'light' | 'dark' | 'system' 
