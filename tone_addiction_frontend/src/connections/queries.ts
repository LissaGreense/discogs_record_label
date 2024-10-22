import { gql } from '@apollo/client';

export const GET_UNIQUE_NAMES = gql`
  query GetUniqueNames {
    uniqueArtists {
      name
    }
    uniqueGenres {
      name
    }
    uniqueStyles {
      name
    }
  }
`;

export const GET_GROUP_WHEN_ARTIST = gql`
  query ReleaseCounts($artist: String!) {
    releaseCounts(artist: $artist) {
      releaseCount
      styleCounts {
        name
        count
      }
      genreCounts {
        name
        count
      }
    }
  }
`;

export const GET_GROUP_WHEN_GENRE = gql`
  query ReleaseCounts($genre: String!) {
    releaseCounts(genre: $genre) {
      releaseCount
      styleCounts {
        name
        count
      }
      artistCounts {
        name
        count
      }
    }
  }
`;

export const GET_GROUP_WHEN_STYLE = gql`
  query ReleaseCounts($style: String!){
    releaseCounts(style: $style) {
      releaseCount
      artistCounts {
        name
        count
      }
      genreCounts {
        name
        count
      }
    }
  }
`;