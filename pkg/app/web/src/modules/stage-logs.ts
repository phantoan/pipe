import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import { LogBlock as LogBlockModel } from "pipe/pkg/app/web/model/logblock_pb";
import { getStageLog } from "../api/stage-log";

export { LogSeverity } from "pipe/pkg/app/web/model/logblock_pb";

export type LogBlock = LogBlockModel.AsObject;

export type StageLog = {
  deploymentId: string;
  stageId: string;
  logBlocks: LogBlock[];
  completed: boolean;
};

// NOTE: Use deploymentId + stageId as record key.
type StageLogs = Record<string, StageLog>;
const initialState: StageLogs = {};

export const createActiveStageKey = (props: {
  deploymentId: string;
  stageId: string;
}): string => `${props.deploymentId}${props.stageId}`;

export const fetchStageLog = createAsyncThunk<
  StageLog,
  {
    deploymentId: string;
    stageId: string;
    offsetIndex: number;
    retriedCount: number;
  }
>(
  "stage-logs/fetch",
  async ({ deploymentId, offsetIndex, retriedCount, stageId }) => {
    const response = await getStageLog({
      deploymentId,
      offsetIndex,
      retriedCount,
      stageId,
    });

    return {
      stageId,
      deploymentId,
      logBlocks: response.blocksList,
      completed: response.completed,
    };
  }
);

export const stageLogsSlice = createSlice({
  name: "stageLogs",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchStageLog.pending, (state, action) => {
        const id = createActiveStageKey(action.meta.arg);
        if (state[id]) {
          state[id].completed = false;
        } else {
          state[id] = {
            stageId: action.meta.arg.stageId,
            deploymentId: action.meta.arg.deploymentId,
            logBlocks: [],
            completed: false,
          };
        }
      })
      .addCase(fetchStageLog.fulfilled, (state, action) => {
        const id = createActiveStageKey(action.meta.arg);
        state[id] = action.payload;
        state[id].completed = true;
      })
      .addCase(fetchStageLog.rejected, (state, action) => {
        const id = createActiveStageKey(action.meta.arg);
        state[id].completed = true;
      });
  },
});

export const selectStageLogById = (
  state: StageLogs,
  props: {
    deploymentId: string;
    stageId: string;
  }
): StageLog | null => {
  return state[createActiveStageKey(props)];
};
