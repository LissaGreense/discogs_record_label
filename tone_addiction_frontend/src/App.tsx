import './App.css'
import {Grid } from "@chakra-ui/react";
import { ArtistGroup, GenreGroup, Group, StyleGroup } from "./types/tableGroups.ts";
import { SearchStack } from "./components/SearchStack.tsx";
import { useState, useEffect } from "react";
import { SectionBox } from "./components/SectionBox.tsx";
import {
  fetchGroupWhenArtist,
  fetchGroupWhenGenre,
  fetchGroupWhenStyle
} from "./connections/fetchGroups.ts";
import {FilteredDataTable} from "./components/FiltredDataTable.tsx";

function App() {
  const [selectedGroup, setSelectedGroup] = useState<Group | null>(null);
  const [selectedValue, setSelectedValue] = useState<string>("");
  const [data, setData] = useState<{ artists: ArtistGroup[], genres: GenreGroup[], styles: StyleGroup[] }>({
    artists: [],
    genres: [],
    styles: []
  });

  useEffect(() => {
    const fetchData = async () => {
      if (selectedGroup && selectedValue) {
        try {
          let fetchedData;

          switch (selectedGroup) {
            case Group.Artist:
              fetchedData = await fetchGroupWhenArtist(selectedValue);
              break;
            case Group.Genre:
              fetchedData = await fetchGroupWhenGenre(selectedValue);
              break;
            case Group.Style:
              fetchedData = await fetchGroupWhenStyle(selectedValue);
              break;
            default:
              fetchedData = null;
              break;
          }

          setData({
            artists: fetchedData?.artistCounts  || [],
            genres: fetchedData?.genreCounts    || [],
            styles: fetchedData?.styleCounts    || [],
          });
        } catch (error) {
          console.error("Error fetching group data:", error);
        }
      }
    };

    fetchData();
  }, [selectedGroup, selectedValue]);

  return (
      <>
        <SectionBox title="Search" />
        <SearchStack
            selectedGroup={selectedGroup}
            setSelectedGroup={setSelectedGroup}
            setSelectedValue={setSelectedValue}
        />

        <SectionBox title="Results" />
        <Grid templateColumns="repeat(2, 1fr)" gap={6}>
          <FilteredDataTable selectedGroup={selectedGroup} data={data} />
        </Grid>
      </>
  );
}

export default App;
