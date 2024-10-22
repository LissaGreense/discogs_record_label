import { Box, Divider, AbsoluteCenter, Heading } from "@chakra-ui/react";

type SectionBoxProps = {
  title: string;
};

export const SectionBox = ({ title }: SectionBoxProps) => {
  return (
      <Box position="relative" padding="30">
        <Divider />
        <AbsoluteCenter px="4">
          <Heading as="h3" size="lg">
            {title}
          </Heading>
        </AbsoluteCenter>
      </Box>
  );
};