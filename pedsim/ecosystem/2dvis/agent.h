//
// pedsim - A microscopic pedestrian simulation system.
// Copyright (c) by Christian Gloor
//

#ifndef AGENT_H
#define AGENT_H

#include "item.h"

#include <QGraphicsItem>

class Agent : public Item {

public:
    Agent();

    QRectF boundingRect() const;
    QPainterPath shape() const;
    void paint(QPainter *painter, const QStyleOptionGraphicsItem *option, QWidget *widget);

protected:
    void advance(int step);

private:
    QColor color;
};

#endif
