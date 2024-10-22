import {ChangeEvent, Dispatch, ReactElement, useEffect, useState} from "react";
import {Box, Button, InputGroup, InputLeftAddon, Select, Stack} from "@chakra-ui/react";
import { Group } from "../types/tableGroups.ts";
import { fetchUniqueNames } from "../connections/fetchUniqueNames.ts";

type SearchStackProps = {
  selectedGroup: Group | null;
  setSelectedGroup: (group: Group | null) => void;
  setSelectedValue: (value: string) => void;
};

export const SearchStack = ({ selectedGroup, setSelectedGroup, setSelectedValue }: SearchStackProps): ReactElement => {
  const [artists, setArtists] = useState<string[]>([]);
  const [genres, setGenres] = useState<string[]>([]);
  const [styles, setStyles] = useState<string[]>([]);
  const [artistValue, setArtistValue] = useState<string>("");
  const [styleValue, setStyleValue] = useState<string>("");
  const [genreValue, setGenreValue] = useState<string>("");

  useEffect(() => {
    const getData = async () => {
      try {
        const { artists, genres, styles } = await fetchUniqueNames();
        setArtists(artists);
        setGenres(genres);
        setStyles(styles);
      } catch (error) {
        console.error("Failed to fetch unique names", error);
      }
    };

    getData();
  }, []);

  useEffect(() => {
    if (selectedGroup === null ) {
      setSelectedValue("")
    }
  }, [selectedGroup]);

  const handleSelectChange = (group: Group, setValue: Dispatch<string>) => (e: ChangeEvent<HTMLSelectElement>) => {
    const selectedValue = e.target.value;
    setValue(selectedValue);

    if (selectedValue === "") {
      setSelectedGroup(null);
      setSelectedValue("");
    } else {
      setSelectedGroup(group);
      setSelectedValue(selectedValue);
    }
  };

  const handleReset = () => {
    setSelectedGroup(null);
    setSelectedValue("");
    setArtistValue("");
    setStyleValue("");
    setGenreValue("");
  };

  return (
      <>
      <Stack spacing={4}>
        <InputGroup>
          <InputLeftAddon>Artist</InputLeftAddon>
          <Select
              placeholder="Select an Artist"
              isDisabled={selectedGroup !== null && selectedGroup !== Group.Artist}
              onChange={handleSelectChange(Group.Artist, setArtistValue)}
              value={artistValue}
          >styleName
            {artists.map((artist, index) => (
                <option key={index} value={artist}>
                  {artist}
                </option>
            ))}
          </Select>
        </InputGroup>

        <InputGroup>
          <InputLeftAddon>Style</InputLeftAddon>
          <Select
              placeholder="Select a Style"
              isDisabled={selectedGroup !== null && selectedGroup !== Group.Style}
              onChange={handleSelectChange(Group.Style, setStyleValue)}
              value={styleValue}

          >
            {styles.map((style, index) => (
                <option key={index} value={style}>
                  {style}
                </option>
            ))}
          </Select>
        </InputGroup>

        <InputGroup>
          <InputLeftAddon>Genre</InputLeftAddon>
          <Select
              placeholder="Select a Genre"
              isDisabled={selectedGroup !== null && selectedGroup !== Group.Genre}
              onChange={handleSelectChange(Group.Genre, setGenreValue)}
              value={genreValue}
          >
            {genres.map((genre, index) => (
                <option key={index} value={genre}>
                  {genre}
                </option>
            ))}
          </Select>
        </InputGroup>
      </Stack>
      <Box textAlign="center" mt={4}>
        <Button colorScheme="blue" onClick={handleReset}>
          Reset Selection
        </Button>
      </Box>
  </>

);
};
