import {Tab, Tabs} from "react-bootstrap";
import SharesTab from "./SharesTab";
import {useEffect, useState} from "react";
import Icon from "@mdi/react";
import {mdiConsoleNetwork} from "@mdi/js";
import {getEnvironmentDetail} from "../../../api/metadata";
import DetailTab from "./DetailTab";
import ActionsTab from "./ActionsTab";

const EnvironmentDetail = (props) => {
    const [detail, setDetail] = useState({});

    useEffect(() => {
        getEnvironmentDetail(props.selection.id)
            .then(resp => {
                setDetail(resp.data);
            });
    }, [props.selection]);

    if(detail.environment) {
        return (
            <div>
                <h2><Icon path={mdiConsoleNetwork} size={2} />{" "}{detail.environment.description}</h2>
                <Tabs defaultActiveKey={"shares"} className={"mb-3"}>
                    <Tab eventKey={"shares"} title={"Shares"}>
                        <SharesTab selection={props.selection} />
                    </Tab>
                    <Tab eventKey={"detail"} title={"Detail"}>
                        <DetailTab environment={detail.environment} />
                    </Tab>
                    <Tab eventKey={"actions"} title={"Actions"}>
                        <ActionsTab environment={detail.environment} />
                    </Tab>
                </Tabs>
            </div>
        );
    }
    return <></>;
};

export default EnvironmentDetail;