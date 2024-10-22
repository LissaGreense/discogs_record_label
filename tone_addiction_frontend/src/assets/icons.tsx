import * as React from "react";

// By: ph
// See: https://v0.app/icon/ph/music-note
// Example: <IconPhMusicNote width="24px" height="24px" style={{color: "#000000"}} />

export const IconPhMusicNote = ({
                                  height = "1em",
                                  fill = "currentColor",
                                  focusable = "false",
                                  ...props
                                }: Omit<React.SVGProps<SVGSVGElement>, "children">) => (
    <svg
        role="img"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 256 256"
        height={height}
        focusable={focusable}
        {...props}
    >
      <path
          fill={fill}
          d="m210.3 56.34l-80-24A8 8 0 0 0 120 40v108.26A48 48 0 1 0 136 184V98.75l69.7 20.91A8 8 0 0 0 216 112V64a8 8 0 0 0-5.7-7.66M88 216a32 32 0 1 1 32-32a32 32 0 0 1-32 32m112-114.75l-64-19.2v-31.3L200 70Z"
      />
    </svg>
);
