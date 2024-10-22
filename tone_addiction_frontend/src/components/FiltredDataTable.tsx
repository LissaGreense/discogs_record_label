import {ReactElement} from 'react';
import {GridItem} from "@chakra-ui/react";
import {DataTable} from "./DataTable.tsx";
import {ArtistGroup, GenreGroup, Group, StyleGroup} from "../types/tableGroups.ts";
import {
  artistGroupColumns,
  genreGroupColumns,
  styleGroupColumns
} from "../utils/columnDefinitions.ts";

type FilteredDataTableProps = {
  selectedGroup: Group | null;
  data: {
    artists: ArtistGroup[];
    genres: GenreGroup[];
    styles: StyleGroup[];
  };
};

export const FilteredDataTable = ({
                                    selectedGroup,
                                    data,
                                  }: FilteredDataTableProps): ReactElement => {
  if (selectedGroup === Group.Artist) {
    return (
        <>
          <GridItem>
            <DataTable data={data.genres} columns={genreGroupColumns}/>
          </GridItem>
          <GridItem>
            <DataTable data={data.styles} columns={styleGroupColumns}/>
          </GridItem>
        </>
    );
  } else if (selectedGroup === Group.Genre) {
    return (
        <>
          <GridItem>
            <DataTable data={data.artists} columns={artistGroupColumns}/>
          </GridItem>
          <GridItem>
            <DataTable data={data.styles} columns={styleGroupColumns}/>
          </GridItem>
        </>
    );
  } else if (selectedGroup === Group.Style) {
    return (
        <>
          <GridItem>
            <DataTable data={data.artists} columns={artistGroupColumns}/>
          </GridItem>
          <GridItem>
            <DataTable data={data.genres} columns={genreGroupColumns}/>
          </GridItem>
        </>
    );
  }

  return <></>;
};
