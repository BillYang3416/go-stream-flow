export interface User {
  name: string;
  lineId: string;
  contact: string;
}

export interface Host extends User {}

export interface Player extends User {}
