export interface ActionLog {
    playerId: string;
    message: string;
}

export interface ActionLogResponse {
    event: 'actionlog';
    data: ActionLog;
}