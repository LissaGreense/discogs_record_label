import { GET_UNIQUE_NAMES } from './queries';
import client from "./graphqlClient.ts";
import {Artist, Genre, Style} from "./commonTypes.ts";

type UniqueNamesData = {
  uniqueArtists: Artist[];
  uniqueGenres: Genre[];
  uniqueStyles: Style[];
}

export const fetchUniqueNames = async () => {
  try {
    const { data } = await client.query<UniqueNamesData>({
      query: GET_UNIQUE_NAMES,
    });

    const artists = data.uniqueArtists.map(artist => artist.name);
    const genres = data.uniqueGenres.map(genre => genre.name);
    const styles = data.uniqueStyles.map(style => style.name);

    return { artists, genres, styles };
  } catch (error) {
    console.error("Error fetching unique names:", error);
    throw error;
  }
};