
export type StyleGroup = {
  name: string;
  count: number;
};

export type GenreGroup = {
  name: string;
  count: number;
};

export type ArtistGroup = {
  name: string;
  count: number;
};

export enum Group {
  _,
  Artist,
  Genre,
  Style,
}