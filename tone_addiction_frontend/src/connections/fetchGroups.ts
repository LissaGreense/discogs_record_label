import {GET_GROUP_WHEN_ARTIST, GET_GROUP_WHEN_GENRE, GET_GROUP_WHEN_STYLE} from './queries';
import client from "./graphqlClient.ts";

type GroupData = {
  releaseCount: number;
  styleCounts:  { name: string; count: number }[];
  genreCounts:  { name: string; count: number }[];
  artistCounts: { name: string; count: number }[];
};

const fetchGroupData = async (query: any, variables: Record<string, any>, emptyField: string): Promise<GroupData> => {
  try {
    const { data } = await client.query({
      query,
      variables,
    });

    const releaseCount = data.releaseCounts.releaseCount;
    const styleCounts = data.releaseCounts.styleCounts || [];

    const genreCounts = data.releaseCounts.genreCounts || [];

    const artistCounts = data.releaseCounts.artistCounts || [];

    return {
      releaseCount,
      styleCounts: emptyField === 'style' ? [] : styleCounts,
      genreCounts: emptyField === 'genre' ? [] : genreCounts,
      artistCounts: emptyField === 'artist' ? [] : artistCounts,
    };
  } catch (error) {
    console.error("Error fetching group data:", error);
    throw error;
  }
};

export const fetchGroupWhenArtist = async (artist: string): Promise<GroupData> => {
  return fetchGroupData(GET_GROUP_WHEN_ARTIST, { artist }, 'artist');
};

export const fetchGroupWhenGenre = async (genre: string): Promise<GroupData> => {
  return fetchGroupData(GET_GROUP_WHEN_GENRE, { genre }, 'genre');
};

export const fetchGroupWhenStyle = async (style: string): Promise<GroupData> => {
  return fetchGroupData(GET_GROUP_WHEN_STYLE, { style }, 'style');
};
