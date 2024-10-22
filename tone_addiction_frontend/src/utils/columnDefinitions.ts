import {createColumnHelper} from "@tanstack/react-table";
import {ArtistGroup, GenreGroup, StyleGroup} from "../types/tableGroups.ts";

const artistGroupColumnHelper = createColumnHelper<ArtistGroup>();

export const artistGroupColumns = [
  artistGroupColumnHelper.accessor("name", {
    cell: (info) => info.getValue(),
    header: "Artist"
  }),
  artistGroupColumnHelper.accessor("count", {
    cell: (info) => info.getValue(),
    header: "Releases Count"
  })
];

const genreGroupColumnHelper = createColumnHelper<GenreGroup>();

export const genreGroupColumns = [
  genreGroupColumnHelper.accessor("name", {
    cell: (info) => info.getValue(),
    header: "Genre"
  }),
  genreGroupColumnHelper.accessor("count", {
    cell: (info) => info.getValue(),
    header: "Releases Count"
  })
];


const styleGroupColumnHelper = createColumnHelper<StyleGroup>();

export const styleGroupColumns = [
  styleGroupColumnHelper.accessor("name", {
    cell: (info) => info.getValue(),
    header: "Style"
  }),
  styleGroupColumnHelper.accessor("count", {
    cell: (info) => info.getValue(),
    header: "Releases Count"
  })
];